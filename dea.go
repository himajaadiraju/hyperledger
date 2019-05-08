
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strconv"
    "time"

    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
)

type BaxterChaincode struct {
}

type order struct {
    ObjectType string `json:"docType"`
    SID string `json:"sid"`
    Qty string `json:"qty"`
    BillToAddress string `json:"billingaddress"`
    ShipToAddress string `json:"shippingaddress"`
    Carrier string `json:"carrier"`
    LogisticsResponse string `json:"lresponse"`
    BaxterResponse string `json:"bresponse"`
    PromisedDeliveryDate string `json:"pdod"`
    PromisedShippingDate string `json:"pdos"`
    DriverResponse string `json:"dresponse"`
    ActualShipDate string `json:"ados"`
    EVAD string `json:"evad"`
    ActualDeliveryDate string `json:"adod"`
    BaxterStatus string `json:"bstatus"`
    CarrierStatus string `json:"cstatus"`
    CustomerDeliveryStatus string `json:"customerstatus"`
    CustomerResponse string `json:"cresponse"`
}

type BaxterAdmin struct {
    ObjectType string `json:"docType"`
    BaxterStatus string `json:"bstatus"`
    SID string `json:"sid"`
    Carrier string `json:"carrier"`
    BaxterResponse string `json:"bresponse"`
}

type Logistics struct {
    ObjectType string `json:"docType"`
    BaxterStatus string `json:"bstatus"`
    SID string `json:"sid"`
    PromisedDeliveryDate string `json:"pdod"`
    PromisedShippingDate string `json:"pdos"`
    LogisticsResponse string `json:"lresponse"`
}

type Driver struct {
    ObjectType string `json:"Driver"`
    BaxterStatus string `json:"bstatus"`
    SID string `json:"sid"`
    DriverResponse string `json:"dresponse"`
    ActualShipDate string `json:"ados"`
    EVAD string `json:"evad"`
    ActualDeliveryDate string `json:"adod"`
    CarrierStatus string `json:"cstatus"`
}

type Customer struct {
    ObjectType string `json:"docType"`
    BaxterStatus string `json:"bstatus"`
    SID string `json:"sid"`
    CustomerDeliveryStatus string `json:"customerstatus"`
    CustomerResponse string `json:"cresponse"`
}

type id struct {
    ObjectType string `json:"docType"`
    SID string `json:"sid"`
}

func main() {
    err := shim.Start(new(BaxterChaincode))
    if err != nil {
        fmt.Printf("Error starting Baxter chaincode: %s", err)
    }

}

func (t *BaxterChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

    return shim.Success(nil)

}

func (t *BaxterChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
    function, args := stub.GetFunctionAndParameters()
    fmt.Println("invoke is running " + function)

    if function == "createOrder" {
        return t.createOrder(stub, args)
    } else if function == "transferToBaxterAdmin" {
        return t.transferToBaxterAdmin(stub, args)
    } else if function == "transferToLogistics" {
        return t.transferToLogistics(stub, args)
    } else if function == "transferToDriver" {
        return t.transferToDriver(stub, args)
    } else if function == "readOrder" {
        return t.readOrder(stub, args)
    } else if function == "getHistoryForOrderwithstatus" {
        return t.getHistoryForOrder(stub, args)
    } else if function == "queryProductByOwner" {
        return t.queryProductByOwner(stub, args)
    } else if function == "queryProduct" {
        return t.queryProduct(stub, args)
    } else if function == "initBaxterAdmin" {
        return t.initBaxterAdmin(stub, args)
    } else if function == "initLogistics" {
        return t.initLogistics(stub, args)
    } else if function == "initDriver" {
        return t.initDriver(stub, args)
    } else if function == "initCustomer" {
        return t.initCustomer(stub, args)
    } else if function == "CustomerIssue" {
        return t.CustomerIssue(stub, args)
    } else if function == "readid" {
        return t.readid(stub, args)
    }

    fmt.Println("invoke did not find func: " + function)

    return shim.Error("Received unknown function invocation")

}

func (t *BaxterChaincode) readid(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    var jsonResp string

    kashyap := "kashyap"
    valAsbytes, err := stub.GetState(kashyap)
    if err != nil {

        jsonResp = "{\"Error\":\"Failed to get state for }"

        return shim.Error(jsonResp)

    } else if valAsbytes == nil {

        jsonResp = "{\"Error\":\"Product does not exist: }"

        return shim.Error(jsonResp)

    }
    fmt.Println("- end read Purchase Quotes (success)")
    return shim.Success(valAsbytes)

}

func (t *BaxterChaincode) createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var err error

    if len(args) != 5 {

        return shim.Error("Incorrect number of arguments. Expecting 5")

    }

    // ==== Input sanitation ====

    fmt.Println("start Create Order")

    if len(args[0]) == 7 {
        return shim.Error("1st argument must be a non-empty string and must be of 8 digits")
    }
    if len(args[1]) <= 0 {
        return shim.Error("2nd argument must be a non-empty string")
    }
    if len(args[2]) <= 0 {
        return shim.Error("3rd argument must be a non-empty string")
    }
    if len(args[3]) <= 0 {
        return shim.Error("4th argument must be a non-empty string")
    }
    if len(args[4]) <= 0 {
        return shim.Error("5th argument must be a non-empty string")
    }

    sid := args[0]
    qty := args[1]
    billingaddress := args[2]
    shippingaddress := args[3]
    carrier := "none"
    lresponse := "none"
    bresponse := "none"
    pdod := "1/1/1999"
    pdos := "1/1/1999"
    dresponse := "none"
    ados := "none"
    evad := "none"
    adod := "none"
    bstatus := args[4]
    cstatus := "none"
    customerstatus := "none"
    cresponse := "none"

    // ==== Check if Order ID already exists ====

    sidAsBytes, err := stub.GetState(sid)

    if err != nil {

        return shim.Error("Failed to get product: " + err.Error())

    } else if sidAsBytes != nil {

        fmt.Println("This SID already exists: " + sid)

        return shim.Error("This product already exists: " + sid)

    }

    // ==== Create product object and marshal to JSON ====

    objectType := "order"

    order := &order{objectType, sid, qty, billingaddress, shippingaddress, carrier, lresponse, bresponse, pdod, pdos, dresponse, ados, evad, adod, bstatus, cstatus, customerstatus, cresponse}

    orderJSONasBytes, err := json.Marshal(order)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(sid, orderJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }
    kashyap := "kashyap"
    docType := "id"

    pqtotransfer := id{}
    pqtotransfer.ObjectType = docType
    pqtotransfer.SID = args[0]
    pqJSONasBytes, _ := json.Marshal(pqtotransfer)

    err = stub.PutState(kashyap, pqJSONasBytes) //rewrite the Product

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end init Product")

    return shim.Success(nil)

}

func (t *BaxterChaincode) initBaxterAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var err error

    // ==== Input sanitation ====

    fmt.Println("start Create Order")

    bstatus := "25"
    sid := "none"
    carrier := "none"
    bresponse := "none"

    // ==== Check if Order ID already exists ====

    bstatusAsBytes, err := stub.GetState(bstatus)

    if err != nil {

        return shim.Error("Failed to get BaxterAdmin: " + err.Error())

    } else if bstatusAsBytes != nil {

        fmt.Println("This BaxterAdmin already exists: " + bstatus)

        return shim.Error("This BaxterAdmin already exists: " + bstatus)

    }

    // ==== Create product object and marshal to JSON ====

    objectType := "BaxterAdmin"

    badmin := &BaxterAdmin{objectType, bstatus, sid, carrier, bresponse}

    bstatusJSONasBytes, err := json.Marshal(badmin)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bstatus, bstatusJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end init BaxterAdmin")

    return shim.Success(nil)

}

func (t *BaxterChaincode) initLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var err error

    // ==== Input sanitation ====

    fmt.Println("start Create Order")

    bstatus := "30"
    sid := "none"
    pdod := "none"
    pdos := "none"
    lresponse := "none"

    // ==== Create product object and marshal to JSON ====

    objectType := "Logistics"

    logistics := &Logistics{objectType, bstatus, sid, pdod, pdos, lresponse}

    logisticsJSONasBytes, err := json.Marshal(logistics)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bstatus, logisticsJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end init logistics")

    return shim.Success(nil)

}

func (t *BaxterChaincode) initDriver(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // ==== Input sanitation ====

    fmt.Println("start Create Order")

    bstatus := args[0]
    sid := "none"
    dresponse := "none"
    ados := "none"
    evad := "none"
    adod := "none"
    cstatus := "none"

    // ==== Check if Order ID already exists ====

    bstatusAsBytes, err := stub.GetState(bstatus)

    if err != nil {

        return shim.Error("Failed to get Driver: " + err.Error())

    } else if bstatusAsBytes != nil {

        fmt.Println("This Driver already exists: " + bstatus)

        return shim.Error("This Driver already exists: " + bstatus)

    }

    // ==== Create product object and marshal to JSON ====

    objectType := "Driver"

    driver := &Driver{objectType, bstatus, sid, dresponse, ados, evad, adod, cstatus}

    driverJSONasBytes, err := json.Marshal(driver)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bstatus, driverJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end init driver")

    return shim.Success(nil)

}

func (t *BaxterChaincode) initCustomer(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // ==== Input sanitation ====

    fmt.Println("start Create Order")

    bstatus := args[0]
    sid := "none"
    customerstatus := "none"
    cresponse := "none"

    // ==== Check if Order ID already exists ====

    bstatusAsBytes, err := stub.GetState(bstatus)

    if err != nil {

        return shim.Error("Failed to get Customer: " + err.Error())

    } else if bstatusAsBytes != nil {

        fmt.Println("This Customer already exists: " + bstatus)

        return shim.Error("This Customer already exists: " + bstatus)

    }

    // ==== Create product object and marshal to JSON ====

    objectType := "Customer"

    customer := &Customer{objectType, bstatus, sid, customerstatus, cresponse}

    customerJSONasBytes, err := json.Marshal(customer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bstatus, customerJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end init customer")

    return shim.Success(nil)

}

func (t *BaxterChaincode) transferToBaxterAdmin(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 3 {

        return shim.Error("Incorrect number of arguments. Expecting 3")

    }

    // ==== Input sanitation ====

    fmt.Println("start transfer Order")

    if len(args[0]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    } else if len(args[1]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    } else if len(args[2]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    }

    bstatus := "25"

    bresponse := args[0]
    sid := args[1]
    carrier := args[2]

    // ==== Check if Order ID already exists ====

    sidAsBytes, err := stub.GetState(sid)

    if err != nil {

        return shim.Error("Failed to get sid: " + err.Error())

    }

    orderToTransfer := order{}
    badminToTransfer := BaxterAdmin{}

    err = json.Unmarshal(sidAsBytes, &orderToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    if orderToTransfer.BaxterStatus != "10" {
        return shim.Error("Wrong previous Baxter Status")
    }

    if bresponse == "Approved" {

        orderToTransfer.Carrier = carrier
        orderToTransfer.BaxterResponse = bresponse
        orderToTransfer.BaxterStatus = bstatus

    } else {
        orderToTransfer.BaxterResponse = bresponse
    }

    orderJSONasBytes, err := json.Marshal(orderToTransfer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(sid, orderJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    badminasBytes, err := stub.GetState(bstatus)

    if err != nil {

        return shim.Error("Failed to get bstatus: " + err.Error())

    }
    err = json.Unmarshal(badminasBytes, &badminToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    badminToTransfer.Carrier = carrier
    badminToTransfer.BaxterResponse = bresponse
    badminToTransfer.SID = sid

    badminJSONasBytes, err := json.Marshal(badminToTransfer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bstatus, badminJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end transfer Order")

    return shim.Success(nil)

}

func (t *BaxterChaincode) transferToLogistics(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 4 {

        return shim.Error("Incorrect number of arguments. Expecting 4")

    }

    // ==== Input sanitation ====

    fmt.Println("start transfer Logistics")

    if len(args[0]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    } else if len(args[1]) <= 0 {
        return shim.Error("2nd argument must be a non-empty string")
    } else if len(args[2]) <= 0 {
        return shim.Error("3rd argument must be a non-empty string")
    } else if len(args[3]) <= 0 {
        return shim.Error("4th argument must be a non-empty string")
    }

    bstatus := "30"

    lresponse := args[0]
    sid := args[1]
    pdos := args[2]
    pdod := args[3]

    // ==== Check if Order ID already exists ====

    sidAsBytes, err := stub.GetState(sid)

    if err != nil {

        return shim.Error("Failed to get sid: " + err.Error())

    }

    orderToTransfer := order{}
    logisticsToTransfer := Logistics{}

    err = json.Unmarshal(sidAsBytes, &orderToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    if orderToTransfer.BaxterStatus != "25" {
        return shim.Error("Wrong previous Baxter Status")
    }

    if lresponse == "Approved" {

        orderToTransfer.PromisedDeliveryDate = pdod
        orderToTransfer.PromisedShippingDate = pdos
        orderToTransfer.LogisticsResponse = lresponse
        orderToTransfer.BaxterStatus = bstatus

    } else {
        orderToTransfer.LogisticsResponse = lresponse
    }

    orderJSONasBytes, err := json.Marshal(orderToTransfer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(sid, orderJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    logisticsasBytes, err := stub.GetState(bstatus)

    if err != nil {

        return shim.Error("Failed to get bstatus: " + err.Error())

    }
    err = json.Unmarshal(logisticsasBytes, &logisticsToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    logisticsToTransfer.PromisedDeliveryDate = pdod
    logisticsToTransfer.PromisedShippingDate = pdos
    logisticsToTransfer.SID = sid
    logisticsToTransfer.LogisticsResponse = lresponse

    logisticsJSONasBytes, err := json.Marshal(logisticsToTransfer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bstatus, logisticsJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end transfer logistics")

    return shim.Success(nil)

}

func (t *BaxterChaincode) transferToDriver(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 3 {

        return shim.Error("Incorrect number of arguments. Expecting 3")

    }

    //==== Input sanitation ====

    sid := args[0]

    // ==== Check if Order ID already exists ====

    sidAsBytes, err := stub.GetState(sid)

    if err != nil {

        return shim.Error("Failed to get sid: " + err.Error())

    }

    orderToTransfer := order{}
    driverToTransfer := Driver{}

    err = json.Unmarshal(sidAsBytes, &orderToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    if orderToTransfer.BaxterStatus == "30" {

        ados := args[1]
        evad := args[2]
        bstatus := "X"
        dresponse := "Accepted Order"

        orderToTransfer.ActualShipDate = ados
        orderToTransfer.EVAD = evad
        orderToTransfer.BaxterStatus = bstatus

        orderJSONasBytes, err := json.Marshal(orderToTransfer)

        if err != nil {

            return shim.Error(err.Error())

        }

        err = stub.PutState(sid, orderJSONasBytes)

        if err != nil {

            return shim.Error(err.Error())

        }

        logisticsasBytes, err := stub.GetState(bstatus)

        if err != nil {

            return shim.Error("Failed to get bstatus: " + err.Error())

        }
        err = json.Unmarshal(logisticsasBytes, &driverToTransfer) //unmarshal it aka JSON.parse()

        if err != nil {

            return shim.Error(err.Error())

        }

        driverToTransfer.ActualDeliveryDate = ados
        driverToTransfer.EVAD = evad
        driverToTransfer.SID = sid
        driverToTransfer.BaxterStatus = bstatus
        driverToTransfer.DriverResponse = dresponse

        logisticsJSONasBytes, err := json.Marshal(driverToTransfer)

        if err != nil {

            return shim.Error(err.Error())

        }

        err = stub.PutState(bstatus, logisticsJSONasBytes)

        if err != nil {

            return shim.Error(err.Error())

        }
    } else if orderToTransfer.BaxterStatus == "X" {

        adod := args[1]
        cstatus := args[2]
        bstatus := "85"
        dresponse := "Order Delivered"

        orderToTransfer.ActualDeliveryDate = adod
        orderToTransfer.BaxterStatus = bstatus
        orderToTransfer.DriverResponse = dresponse
        orderToTransfer.CarrierStatus = cstatus

        orderJSONasBytes, err := json.Marshal(orderToTransfer)

        if err != nil {

            return shim.Error(err.Error())

        }

        err = stub.PutState(sid, orderJSONasBytes)

        if err != nil {

            return shim.Error(err.Error())

        }

        logisticsasBytes, err := stub.GetState(bstatus)

        if err != nil {

            return shim.Error("Failed to get bstatus: " + err.Error())

        }
        err = json.Unmarshal(logisticsasBytes, &driverToTransfer) //unmarshal it aka JSON.parse()

        if err != nil {

            return shim.Error(err.Error())

        }

        driverToTransfer.ActualDeliveryDate = adod
        driverToTransfer.BaxterStatus = bstatus
        driverToTransfer.SID = sid
        driverToTransfer.DriverResponse = dresponse
        driverToTransfer.CarrierStatus = cstatus

        logisticsJSONasBytes, err := json.Marshal(driverToTransfer)

        if err != nil {

            return shim.Error(err.Error())

        }

        err = stub.PutState(bstatus, logisticsJSONasBytes)

        if err != nil {

            return shim.Error(err.Error())

        }
    }
    fmt.Println("end transfer logistics")

    return shim.Success(nil)

}

func (t *BaxterChaincode) readOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    var productid, jsonResp string

    var err error

    if len(args) != 1 {

        return shim.Error("Incorrect number of arguments. Expecting Name of the Product to query")

    }

    productid = args[0]

    valAsbytes, err := stub.GetState(productid) //get the medicine from chaincode state

    if err != nil {

        jsonResp = "{\"Error\":\"Failed to get state for " + productid + "\"}"

        return shim.Error(jsonResp)

    } else if valAsbytes == nil {

        jsonResp = "{\"Error\":\"Order does not exist: " + productid + "\"}"

        return shim.Error(jsonResp)

    }

    return shim.Success(valAsbytes)

}

func (t *BaxterChaincode) queryProductByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // 0

    // "bob"

    if len(args) < 1 {

        return shim.Error("Incorrect number of arguments. Expecting 1")

    }

    owner := args[0]

    queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"marble\",\"owner\":\"%s\"}}", owner)

    queryResults, err := getQueryResultForQueryString(stub, queryString)

    if err != nil {

        return shim.Error(err.Error())

    }

    return shim.Success(queryResults)

}

//QueryMedicine

func (t *BaxterChaincode) queryProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // 0

    // "queryString"

    if len(args) < 1 {

        return shim.Error("Incorrect number of arguments. Expecting 1")

    }

    queryString := args[0]

    queryResults, err := getQueryResultForQueryString(stub, queryString)

    if err != nil {

        return shim.Error(err.Error())

    }

    return shim.Success(queryResults)

}

func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

    fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

    resultsIterator, err := stub.GetQueryResult(queryString)

    if err != nil {

        return nil, err

    }

    defer resultsIterator.Close()

    // buffer is a JSON array containing QueryRecords

    var buffer bytes.Buffer

    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false

    for resultsIterator.HasNext() {

        queryResponse, err := resultsIterator.Next()

        if err != nil {

            return nil, err

        }

        // Add a comma before array members, suppress it for the first array member

        if bArrayMemberAlreadyWritten == true {

            buffer.WriteString(",")

        }

        buffer.WriteString("{\"Key\":")

        buffer.WriteString("\"")

        buffer.WriteString(queryResponse.Key)

        buffer.WriteString("\"")

        buffer.WriteString(", \"Record\":")

        // Record is a JSON object, so we write as-is

        buffer.WriteString(string(queryResponse.Value))

        buffer.WriteString("}")

        bArrayMemberAlreadyWritten = true

    }

    buffer.WriteString("]")

    fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

    return buffer.Bytes(), nil

}

func (t *BaxterChaincode) getHistoryForOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) < 1 {

        return shim.Error("Incorrect number of arguments. Expecting 1")

    }

    productid := args[0]

    fmt.Printf("- start getHistoryForProduct: %s\n", productid)

    resultsIterator, err := stub.GetHistoryForKey(productid)

    if err != nil {

        return shim.Error(err.Error())

    }

    defer resultsIterator.Close()

    // buffer is a JSON array containing historic values for the medicine

    var buffer bytes.Buffer

    buffer.WriteString("[")

    bArrayMemberAlreadyWritten := false

    for resultsIterator.HasNext() {

        response, err := resultsIterator.Next()

        if err != nil {

            return shim.Error(err.Error())

        }

        // Add a comma before array members, suppress it for the first array member

        if bArrayMemberAlreadyWritten == true {

            buffer.WriteString(",")

        }

        buffer.WriteString("{\"TxId\":")

        buffer.WriteString("\"")

        buffer.WriteString(response.TxId)

        buffer.WriteString("\"")

        buffer.WriteString(", \"Value\":")

        // if it was a delete operation on given key, then we need to set the

        //corresponding value null. Else, we will write the response.Value

        //as-is (as the Value itself a JSON medicine)

        if response.IsDelete {

            buffer.WriteString("null")

        } else {

            buffer.WriteString(string(response.Value))

        }

        buffer.WriteString(", \"Timestamp\":")

        buffer.WriteString("\"")

        buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())

        buffer.WriteString("\"")

        buffer.WriteString(", \"IsDelete\":")

        buffer.WriteString("\"")

        buffer.WriteString(strconv.FormatBool(response.IsDelete))

        buffer.WriteString("\"")

        buffer.WriteString("}")

        bArrayMemberAlreadyWritten = true

    }

    buffer.WriteString("]")

    fmt.Printf("- getHistoryForProduct returning:\n%s\n", buffer.String())

    return shim.Success(buffer.Bytes())

}

func (t *BaxterChaincode) CustomerIssue(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 3 {

        return shim.Error("Incorrect number of arguments. Expecting 3")

    }

    // ==== Input sanitation ====

    fmt.Println("start transfer Order")

    if len(args[0]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    } else if len(args[1]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    } else if len(args[2]) <= 0 {
        return shim.Error("1st argument must be a non-empty string")
    }

    bfstatus := "99"

    sid := args[0]
    cresponse := args[1]
    customerstatus := args[2]

    // ==== Check if Order ID already exists ====

    sidAsBytes, err := stub.GetState(sid)

    if err != nil {

        return shim.Error("Failed to get sid: " + err.Error())

    }

    orderToTransfer := order{}
    customerToTransfer := Customer{}

    err = json.Unmarshal(sidAsBytes, &orderToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    if orderToTransfer.BaxterStatus != "85" {
        return shim.Error("Wrong previous Baxter Status")
    }

    if cresponse == "Delivery Completed" {

        orderToTransfer.CustomerResponse = cresponse
        orderToTransfer.CustomerDeliveryStatus = customerstatus

    } else {
        orderToTransfer.CustomerResponse = cresponse
        orderToTransfer.CustomerDeliveryStatus = customerstatus
        orderToTransfer.BaxterStatus = bfstatus
    }

    orderJSONasBytes, err := json.Marshal(orderToTransfer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(sid, orderJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    bfadminasBytes, err := stub.GetState(bfstatus)

    if err != nil {

        return shim.Error("Failed to get bfstatus: " + err.Error())

    }
    err = json.Unmarshal(bfadminasBytes, &customerToTransfer) //unmarshal it aka JSON.parse()

    if err != nil {

        return shim.Error(err.Error())

    }

    customerToTransfer.CustomerDeliveryStatus = customerstatus
    customerToTransfer.CustomerResponse = cresponse
    customerToTransfer.SID = sid
    customerToTransfer.BaxterStatus = bfstatus

    bfadminJSONasBytes, err := json.Marshal(customerToTransfer)

    if err != nil {

        return shim.Error(err.Error())

    }

    err = stub.PutState(bfstatus, bfadminJSONasBytes)

    if err != nil {

        return shim.Error(err.Error())

    }

    fmt.Println("end transfer Order")

    return shim.Success(nil)

}