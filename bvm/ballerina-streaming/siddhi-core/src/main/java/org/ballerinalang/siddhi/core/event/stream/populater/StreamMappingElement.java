/*
 * Copyright (c) 2016, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */
package org.ballerinalang.siddhi.core.event.stream.populater;

/**
 * Element to hold information about StreamEvent population from the given data[].
 */
public class StreamMappingElement {
    private int fromPosition;   //position in data[]
    private int[] toPosition;   //position in StateEvent (PreOutput,Output)/data[]

    public int[] getToPosition() {
        return toPosition;
    }

    public void setToPosition(int[] toPosition) {
        this.toPosition = toPosition;
    }

    public int getFromPosition() {
        return fromPosition;
    }

    public void setFromPosition(int fromPosition) {
        this.fromPosition = fromPosition;
    }

}
