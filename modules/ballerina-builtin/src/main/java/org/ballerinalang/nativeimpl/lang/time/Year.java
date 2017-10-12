/*
*   Copyright (c) 2017, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
*/
package org.ballerinalang.nativeimpl.lang.time;

import org.ballerinalang.bre.Context;
import org.ballerinalang.model.types.TypeKind;
import org.ballerinalang.model.values.BInteger;
import org.ballerinalang.model.values.BStruct;
import org.ballerinalang.model.values.BValue;
import org.ballerinalang.natives.annotations.Argument;
import org.ballerinalang.natives.annotations.Attribute;
import org.ballerinalang.natives.annotations.BallerinaAnnotation;
import org.ballerinalang.natives.annotations.BallerinaFunction;
import org.ballerinalang.natives.annotations.ReturnType;

/**
 * Get the year value for the given time.
 *
 * @since 0.89
 */
@BallerinaFunction(
        packageName = "ballerina.lang.time",
        functionName = "year",
        args = {@Argument(name = "time", type = TypeKind.STRUCT, structType = "Time",
                          structPackage = "ballerina.lang.time")},
        returnType = {@ReturnType(type = TypeKind.INT)},
        isPublic = true
)
@BallerinaAnnotation(annotationName = "Description", attributes = { @Attribute(name = "value",
       value = "Get the year value of the given the Time.")})
@BallerinaAnnotation(annotationName = "Return", attributes = {@Attribute(name = "int) ",
       value = "Year of the given time value")})
public class Year extends AbstractTimeFunction {

    @Override
    public BValue[] execute(Context context) {
        BStruct timeStruct = ((BStruct) getRefArgument(context, 0));
        return new BValue[]{new BInteger(getYear(timeStruct))};
    }
}