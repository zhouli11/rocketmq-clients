package org.apache.rocketmq.client.message;

import java.util.Collections;
import java.util.HashSet;
import java.util.Set;

public class MessageConst {
    public static final String PROPERTY_KEYS = "KEYS";
    public static final String PROPERTY_TAGS = "TAGS";
    public static final String PROPERTY_WAIT_STORE_MSG_OK = "WAIT";
    public static final String PROPERTY_DELAY_TIME_LEVEL = "DELAY";
    public static final String PROPERTY_RETRY_TOPIC = "RETRY_TOPIC";
    public static final String PROPERTY_REAL_TOPIC = "REAL_TOPIC";
    /**
     * Full qualified name of the topic, typically with namespace, retry prefix or DLQ prefix.
     */
    public static final String PROPERTY_FQN_TOPIC = "__FQN_TOPIC";

    public static final String PROPERTY_REAL_QUEUE_ID = "REAL_QID";
    public static final String PROPERTY_TRANSACTION_PREPARED = "TRAN_MSG";
    public static final String PROPERTY_PRODUCER_GROUP = "PGROUP";
    public static final String PROPERTY_MIN_OFFSET = "MIN_OFFSET";
    public static final String PROPERTY_MAX_OFFSET = "MAX_OFFSET";
    public static final String PROPERTY_BUYER_ID = "BUYER_ID";
    public static final String PROPERTY_ORIGIN_MESSAGE_ID = "ORIGIN_MESSAGE_ID";
    public static final String PROPERTY_TRANSFER_FLAG = "TRANSFER_FLAG";
    public static final String PROPERTY_CORRECTION_FLAG = "CORRECTION_FLAG";
    public static final String PROPERTY_MQ2_FLAG = "MQ2_FLAG";
    public static final String PROPERTY_RECONSUME_TIME = "RECONSUME_TIME";
    public static final String PROPERTY_MSG_REGION = "MSG_REGION";
    public static final String PROPERTY_TRACE_SWITCH = "TRACE_ON";
    public static final String PROPERTY_UNIQ_CLIENT_MESSAGE_ID_KEYIDX = "UNIQ_KEY";
    public static final String PROPERTY_EXTEND_UNIQ_INFO = "EXTEND_UNIQ_INFO";
    public static final String PROPERTY_MAX_RECONSUME_TIMES = "MAX_RECONSUME_TIMES";
    public static final String PROPERTY_CONSUME_START_TIMESTAMP = "CONSUME_START_TIME";
    public static final String PROPERTY_INSTANCE_ID = "INSTANCE_ID";
    public static final String PROPERTY_TIMER_DELAY_SEC = "TIMER_DELAY_SEC";
    public static final String PROPERTY_TIMER_DELIVER_MS = "TIMER_DELIVER_MS";
    public static final String PROPERTY_TIMER_IN_MS = "TIMER_IN_MS";
    public static final String PROPERTY_TIMER_OUT_MS = "TIMER_OUT_MS";
    public static final String PROPERTY_TIMER_ENQUEUE_MS = "TIMER_ENQUEUE_MS";
    public static final String PROPERTY_TIMER_DEQUEUE_MS = "TIMER_DEQUEUE_MS";
    public static final String PROPERTY_TIMER_ROLL_TIMES = "TIMER_ROLL_TIMES";
    public static final String PROPERTY_TIMER_DEL_UNIQ_KEY = "TIMER_DEL_UNIQKEY";
    public static final String PROPERTY_TIMER_DELAY_LEVEL = "TIMER_DELAY_LEVEL";
    public static final String PROPERTY_POP_CK = "POP_CK";
    public static final String PROPERTY_ACK_HOST_ADDRESS = "__ACK_HOST";
    public static final String PROPERTY_POP_CK_OFFSET = "POP_CK_OFFSET";
    public static final String PROPERTY_FIRST_POP_TIME = "1ST_POP_TIME";
    public static final String PROPERTY_SHARDING_KEY = "__SHARDINGKEY";
    public static final String PROPERTY_VTOA_TUNNEL_ID = "VTOA_TUNNEL_ID";

    public static final String KEY_SEPARATOR = " ";


    private static final Set<String> MSG_SYS_PROPERTY_SET_0 = new HashSet<String>();

    static {
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TRACE_SWITCH);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_MSG_REGION);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_KEYS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TAGS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_WAIT_STORE_MSG_OK);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_DELAY_TIME_LEVEL);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_RETRY_TOPIC);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_REAL_TOPIC);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_FQN_TOPIC);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_REAL_QUEUE_ID);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TRANSACTION_PREPARED);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_PRODUCER_GROUP);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_MIN_OFFSET);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_MAX_OFFSET);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_BUYER_ID);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_ORIGIN_MESSAGE_ID);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TRANSFER_FLAG);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_CORRECTION_FLAG);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_MQ2_FLAG);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_RECONSUME_TIME);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_UNIQ_CLIENT_MESSAGE_ID_KEYIDX);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_MAX_RECONSUME_TIMES);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_CONSUME_START_TIMESTAMP);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TIMER_IN_MS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TIMER_OUT_MS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TIMER_ENQUEUE_MS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TIMER_DEQUEUE_MS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_TIMER_ROLL_TIMES);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_POP_CK);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_ACK_HOST_ADDRESS);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_POP_CK_OFFSET);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_FIRST_POP_TIME);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_EXTEND_UNIQ_INFO);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_INSTANCE_ID);
        MSG_SYS_PROPERTY_SET_0.add(PROPERTY_VTOA_TUNNEL_ID);
    }

    public static final Set<String> MSG_SYS_PROPERTY_SET = Collections.unmodifiableSet(MSG_SYS_PROPERTY_SET_0);

    private MessageConst() {
    }
}
