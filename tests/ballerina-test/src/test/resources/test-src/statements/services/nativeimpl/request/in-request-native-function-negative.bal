import ballerina/net.http;
import ballerina/mime;

function testGetContentLength (http:Request req) returns (string) {
    if (req.hasHeader("content-legth")) {
        return req.getHeader("content-length");
    }
    return "Content-length is not found";
}

function testGetHeader (http:Request req, string key) returns (string) {
    if (req.hasHeader(key)) {
        return req.getHeader(key);
    }
    return "Header not found!";
}

function testGetJsonPayload (http:Request req) returns (json | mime:EntityError) {
    return req.getJsonPayload();
}

function testGetMethod (http:Request req) returns (string) {
    string method = req.method;
    return method;
}

function testGetRequestURL (http:Request req) returns (string) {
    string url = req.rawPath;
    if (url == null || url == "") {
        url = "no url";
    }
    return url;
}

function testGetStringPayload (http:Request req) returns (string | null | mime:EntityError) {
    return req.getStringPayload();
}

function testGetBinaryPayload (http:Request req) returns (blob | mime:EntityError) {
    return req.getBinaryPayload();
}

function testGetXmlPayload (http:Request req) returns (xml | mime:EntityError) {
    return req.getXmlPayload();
}

function testSetHeader (http:Request req, string key, string value) returns (http:Request) {
    req.setHeader(key, value);
    return req;
}

function testGetEntity (http:Request req) returns (mime:Entity | mime:EntityError) {
    return req.getEntity();
}

function testRemoveHeader (http:Request req, string key) returns (http:Request) {
    req.removeHeader(key);
    return req;
}

function testRemoveAllHeaders (http:Request req) returns (http:Request) {
    req.removeAllHeaders();
    return req;
}
