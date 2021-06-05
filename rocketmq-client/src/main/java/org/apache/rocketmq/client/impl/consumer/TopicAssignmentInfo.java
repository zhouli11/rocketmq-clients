package org.apache.rocketmq.client.impl.consumer;

import apache.rocketmq.v1.LoadAssignment;
import java.util.ArrayList;
import java.util.List;
import lombok.EqualsAndHashCode;
import lombok.Getter;
import lombok.ToString;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.RandomUtils;
import org.apache.rocketmq.client.message.MessageQueue;
import org.apache.rocketmq.client.route.Partition;

@Slf4j
@ToString
@EqualsAndHashCode
public class TopicAssignmentInfo {
    private static final ThreadLocal<Integer> PARTITION_INDEX = new ThreadLocal<Integer>();

    @Getter
    private final List<Assignment> assignmentList;

    static {
        PARTITION_INDEX.set(RandomUtils.nextInt());
    }

    public TopicAssignmentInfo(List<LoadAssignment> loadAssignmentList) {
        this.assignmentList = new ArrayList<Assignment>();

        for (LoadAssignment item : loadAssignmentList) {
            MessageQueue messageQueue =
                    new MessageQueue(new Partition(item.getPartition()));

            MessageRequestMode mode = MessageRequestMode.POP;
            switch (item.getMode()) {
                case PULL:
                    mode = MessageRequestMode.PULL;
                    break;
                case POP:
                    mode = MessageRequestMode.POP;
                    break;
                default:
                    log.warn("Unknown message request mode={}, default to pop.", item.getMode());
            }
            assignmentList.add(new Assignment(messageQueue, mode));
        }
    }

    public static int getNextPartitionIndex() {
        Integer index = PARTITION_INDEX.get();
        if (null == index) {
            index = -1;
        }
        index += 1;
        index = Math.abs(index);
        PARTITION_INDEX.set(index);
        return index;
    }


}
