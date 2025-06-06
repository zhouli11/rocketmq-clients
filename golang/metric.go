/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package golang

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"

	"contrib.go.opencensus.io/exporter/ocagent"
	"github.com/apache/rocketmq-clients/golang/v5/pkg/utils"
	v2 "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type InvocationStatus string

const (
	InvocationStatus_SUCCESS InvocationStatus = "success"
	InvocationStatus_FAILURE InvocationStatus = "failure"
)

var (
	topicTag, _            = tag.NewKey("topic")
	clientIdTag, _         = tag.NewKey("client_id")
	invocationStatusTag, _ = tag.NewKey("invocation_status")
	consumerGroupTag, _    = tag.NewKey("consumer_group")

	PublishMLatencyMs         = stats.Int64("publish_latency", "Publish latency in milliseconds", "ms")
	ConsumeDeliveryMLatencyMs = stats.Int64("delivery_latency", "Time spent delivering messages from servers to clients", "ms")
	ConsumeAwaitMLatencyMs    = stats.Int64("await_time", "Client side queuing time of messages before getting processed", "ms")
	ConsumeProcessMLatencyMs  = stats.Int64("process_time", "Process message time", "ms")

	PublishLatencyView = view.View{
		Name:        "rocketmq_send_cost_time",
		Description: "Publish latency",
		Measure:     PublishMLatencyMs,
		Aggregation: view.Distribution(1, 5, 10, 20, 50, 200, 500),
		TagKeys:     []tag.Key{topicTag, clientIdTag, invocationStatusTag},
	}

	ConsumeDeliveryLatencyView = view.View{
		Name:        "rocketmq_delivery_latency",
		Description: "Message delivery latency",
		Measure:     ConsumeDeliveryMLatencyMs,
		Aggregation: view.Distribution(1, 5, 10, 20, 50, 200, 500),
		TagKeys:     []tag.Key{topicTag, clientIdTag, consumerGroupTag},
	}

	ConsumeAwaitTimeView = view.View{
		Name:        "rocketmq_await_time",
		Description: "Message await time",
		Measure:     ConsumeAwaitMLatencyMs,
		Aggregation: view.Distribution(1, 5, 20, 100, 1000, 5000, 10000),
		TagKeys:     []tag.Key{topicTag, clientIdTag, consumerGroupTag},
	}

	ConsumeProcessTimeView = view.View{
		Name:        "rocketmq_process_time",
		Description: "Message process time",
		Measure:     ConsumeProcessMLatencyMs,
		Aggregation: view.Distribution(1, 5, 10, 100, 1000, 10000, 60000),
		TagKeys:     []tag.Key{topicTag, clientIdTag, consumerGroupTag, invocationStatusTag},
	}
)

func init() {
	if err := view.Register(&PublishLatencyView, &ConsumeDeliveryLatencyView, &ConsumeAwaitTimeView, &ConsumeProcessTimeView); err != nil {
		sugarBaseLogger.Fatalf("failed to register views: %v", err)
	}
	view.SetReportingPeriod(time.Minute)
}

type defaultClientMeter struct {
	enabled     atomic.Bool
	endpoints   *v2.Endpoints
	ocaExporter view.Exporter
	mutex       sync.Mutex
}

func (dcm *defaultClientMeter) shutdown() {
	if !dcm.enabled.Load() {
		return
	}
	dcm.mutex.Lock()
	defer dcm.mutex.Unlock()
	view.UnregisterExporter(dcm.ocaExporter)
	if dcm.ocaExporter != nil {
		exporter, ok := dcm.ocaExporter.(*ocagent.Exporter)
		if ok {
			err := exporter.Stop()
			if err != nil {
				sugarBaseLogger.Errorf("ocExporter stop failed, err=%w", err)
			}
		}
	}
}

func (dcm *defaultClientMeter) start() {
	if !dcm.enabled.Load() {
		return
	}
	view.RegisterExporter(dcm.ocaExporter)
}

var NewDefaultClientMeter = func(exporter view.Exporter, on bool, endpoints *v2.Endpoints, clientID string) *defaultClientMeter {
	return &defaultClientMeter{
		enabled:     *atomic.NewBool(on),
		endpoints:   endpoints,
		ocaExporter: exporter,
	}
}

type MessageMeterInterceptor interface {
	MessageInterceptor
}

type defaultMessageMeterInterceptor struct {
	clientMeterProvider ClientMeterProvider
}

type ClientMeterProvider interface {
	Reset(metric *v2.Metric)
	isEnabled() bool
	getClientID() string
	getClientImpl() isClient
}

var _ = ClientMeterProvider(&defaultClientMeterProvider{})

type defaultClientMeterProvider struct {
	client      Client
	clientMeter *defaultClientMeter
	globalMutex sync.Mutex
}

func (dcmp *defaultClientMeterProvider) getClientImpl() isClient {
	if dc, ok := dcmp.client.(*defaultClient); ok {
		return dc.clientImpl
	}
	return nil
}

var _ = MessageMeterInterceptor(&defaultMessageMeterInterceptor{})

var NewDefaultMessageMeterInterceptor = func(clientMeterProvider ClientMeterProvider) *defaultMessageMeterInterceptor {
	return &defaultMessageMeterInterceptor{
		clientMeterProvider: clientMeterProvider,
	}
}

func (dmmi *defaultMessageMeterInterceptor) doBeforeConsumeMessage(messageCommons []*MessageCommon) error {
	if len(messageCommons) == 0 {
		// Should never reach here.
		return nil
	}
	clientImpl := dmmi.clientMeterProvider.getClientImpl()
	if clientImpl == nil {
		return nil
	}
	var pc PushConsumer
	var ok bool
	if pc, ok = clientImpl.(PushConsumer); !ok {
		return nil
	}
	consumerGroup := pc.GetGroupName()
	clientId := dmmi.clientMeterProvider.getClientID()
	if len(consumerGroup) == 0 {
		sugarBaseLogger.Errorf("[Bug] consumerGroup is not recognized, clientId=%s", clientId)
		return nil
	}
	for _, messageCommon := range messageCommons {
		if messageCommon.decodeStopwatch == nil {
			continue
		}
		duration := time.Since(*messageCommon.decodeStopwatch)
		err := stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Insert(topicTag, messageCommon.topic), tag.Insert(clientIdTag, dmmi.clientMeterProvider.getClientID()), tag.Insert(consumerGroupTag, consumerGroup)}, ConsumeAwaitMLatencyMs.M(duration.Milliseconds()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (dmmi *defaultMessageMeterInterceptor) doAfterConsumeMessage(messageCommons []*MessageCommon, duration time.Duration, status MessageHookPointsStatus) error {
	if len(messageCommons) == 0 {
		// Should never reach here.
		return nil
	}
	clientImpl := dmmi.clientMeterProvider.getClientImpl()
	if clientImpl == nil {
		return nil
	}
	var pc PushConsumer
	var ok bool
	if pc, ok = clientImpl.(PushConsumer); !ok {
		return nil
	}
	consumerGroup := pc.GetGroupName()
	clientId := dmmi.clientMeterProvider.getClientID()
	if len(consumerGroup) == 0 {
		sugarBaseLogger.Errorf("[Bug] consumerGroup is not recognized, clientId=%s", clientId)
		return nil
	}

	invocationStatus := InvocationStatus_FAILURE
	if status == MessageHookPointsStatus_OK {
		invocationStatus = InvocationStatus_SUCCESS
	}
	for _, messageCommon := range messageCommons {
		err := stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Insert(topicTag, messageCommon.topic), tag.Insert(clientIdTag, dmmi.clientMeterProvider.getClientID()), tag.Insert(consumerGroupTag, consumerGroup), tag.Insert(invocationStatusTag, string(invocationStatus))}, ConsumeProcessMLatencyMs.M(duration.Milliseconds()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (dmmi *defaultMessageMeterInterceptor) doAfterReceiveMessage(messageCommons []*MessageCommon, duration time.Duration, status MessageHookPointsStatus) error {
	if len(messageCommons) == 0 {
		// Should never reach here.
		return nil
	}
	clientImpl := dmmi.clientMeterProvider.getClientImpl()
	if clientImpl == nil {
		return nil
	}
	var pc PushConsumer
	var ok bool
	if pc, ok = clientImpl.(PushConsumer); !ok {
		return nil
	}
	consumerGroup := pc.GetGroupName()
	clientId := dmmi.clientMeterProvider.getClientID()
	if len(consumerGroup) == 0 {
		sugarBaseLogger.Errorf("[Bug] consumerGroup is not recognized, clientId=%s", clientId)
		return nil
	}

	for _, messageCommon := range messageCommons {
		if messageCommon.deliveryTimestamp == nil {
			continue
		}
		latency := time.Since(*messageCommon.deliveryTimestamp)
		err := stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Insert(topicTag, messageCommon.topic), tag.Insert(clientIdTag, dmmi.clientMeterProvider.getClientID()), tag.Insert(consumerGroupTag, consumerGroup)}, ConsumeDeliveryMLatencyMs.M(latency.Milliseconds()))
		if err != nil {
			return err
		}
	}

	return nil
}

func (dmmi *defaultMessageMeterInterceptor) doBefore(messageHookPoints MessageHookPoints, messageCommons []*MessageCommon) error {
	if !dmmi.clientMeterProvider.isEnabled() {
		return nil
	}
	switch messageHookPoints {
	case MessageHookPoints_CONSUME:
		return dmmi.doBeforeConsumeMessage(messageCommons)
	default:
		break
	}
	return nil
}

func (dmmi *defaultMessageMeterInterceptor) doAfterSendMessage(messageCommons []*MessageCommon, duration time.Duration, status MessageHookPointsStatus) error {
	invocationStatus := InvocationStatus_FAILURE
	if status == MessageHookPointsStatus_OK {
		invocationStatus = InvocationStatus_SUCCESS
	}
	for _, messageCommon := range messageCommons {
		err := stats.RecordWithTags(context.Background(), []tag.Mutator{tag.Insert(topicTag, messageCommon.topic), tag.Insert(clientIdTag, dmmi.clientMeterProvider.getClientID()), tag.Insert(invocationStatusTag, string(invocationStatus))}, PublishMLatencyMs.M(duration.Milliseconds()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (dmmi *defaultMessageMeterInterceptor) doAfter(messageHookPoints MessageHookPoints, messageCommons []*MessageCommon, duration time.Duration, status MessageHookPointsStatus) error {
	if !dmmi.clientMeterProvider.isEnabled() {
		return nil
	}
	switch messageHookPoints {
	case MessageHookPoints_SEND:
		return dmmi.doAfterSendMessage(messageCommons, duration, status)
	case MessageHookPoints_CONSUME:
		return dmmi.doAfterConsumeMessage(messageCommons, duration, status)
	case MessageHookPoints_RECEIVE:
		return dmmi.doAfterReceiveMessage(messageCommons, duration, status)
	default:
		break
	}
	return nil
}
func (dcmp *defaultClientMeterProvider) isEnabled() bool {
	return dcmp.clientMeter.enabled.Load()
}
func (dcmp *defaultClientMeterProvider) getClientID() string {
	return dcmp.client.GetClientID()
}
func (dcmp *defaultClientMeterProvider) Reset(metric *v2.Metric) {
	dcmp.globalMutex.Lock()
	defer dcmp.globalMutex.Unlock()
	endpoints := metric.GetEndpoints()
	if dcmp.clientMeter.enabled.Load() && metric.GetOn() && utils.CompareEndpoints(dcmp.clientMeter.endpoints, endpoints) {
		sugarBaseLogger.Infof("metric settings is satisfied by the current message meter, clientId=%s", dcmp.client.GetClientID())
		return
	}

	if !metric.GetOn() {
		dcmp.clientMeter.shutdown()
		sugarBaseLogger.Infof("metric is off, clientId=%s", dcmp.client.GetClientID())
		dcmp.clientMeter = NewDefaultClientMeter(nil, false, nil, dcmp.client.GetClientID())
		return
	}
	agentAddr := utils.ParseAddress(utils.SelectAnAddress(endpoints))
	exporter, err := ocagent.NewExporter(
		ocagent.WithInsecure(),
		ocagent.WithTLSCredentials(credentials.NewTLS(defaultConnOptions.TLS)),
		ocagent.WithAddress(agentAddr),
		ocagent.WithGRPCDialOption(grpc.WithChainUnaryInterceptor(dcmp.invokeWithSign())),
	)
	if err != nil {
		sugarBaseLogger.Errorf("exception raised when resetting message meter, clientId=%s", dcmp.client.GetClientID())
		return
	}
	// Reset message meter.
	dcmp.clientMeter.shutdown()
	dcmp.clientMeter = NewDefaultClientMeter(exporter, true, endpoints, dcmp.client.GetClientID())
	dcmp.clientMeter.start()
	sugarBaseLogger.Infof("metrics is on, endpoints=%v, clientId=%s", endpoints, dcmp.client.GetClientID())
}

var NewDefaultClientMeterProvider = func(client *defaultClient) ClientMeterProvider {
	cmp := &defaultClientMeterProvider{
		client:      client,
		clientMeter: NewDefaultClientMeter(nil, false, nil, "nil"),
	}
	client.registerMessageInterceptor(NewDefaultMessageMeterInterceptor(cmp))
	return cmp
}

var _ = ClientMeterProvider(&defaultClientMeterProvider{})

func (dcmp *defaultClientMeterProvider) invokeWithSign() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx := dcmp.client.Sign(ctx)
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}
