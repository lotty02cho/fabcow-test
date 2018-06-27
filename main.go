/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

type Owner struct {
	Owner_id         string `json:"Owner_id"`
	Owner_nm         string `json:"Owner_nm"`
	Owner_addr       string `json:"Owner_addr"`
	Livestock        string `json:"Livestock"`
	Owner_user_nm    string `json:"Owner_user_nm"`
	Owner_user_birth string `json:"Owner_user_birth"`
	Remarks          []Remark
}

// Define the cow structure, with 4 properties.  Structure tags are used by encoding/json library
type Cow struct {
	Id_no      string `json:"Id_no"`
	Birth_date string `json:"Birth_date"`
	Sex        string `json:"Sex"`
	Father_id  string `json:"Father_id"`
	Mother_id  string `json:"Mother_id"`
	Origin     string `json:"Origin"`
	Owner      Owner
	Remarks    []Remark
}

type HACCP struct {
	Farm_id       string `json:"Farm_id"`
	Farm_nm       string `json:"Farm_nm"`
	Farm_addr     string `json:"Farm_addr"`
	Apply_item    string `json:"Apply_item"`
	Validity_date string `json:"Validity_date"`
}

type RFID struct {
	Id_no   string `json:"Id_no"`
	Rfid_no string `json:"Rfid_no"`
}

type Bundle struct {
	Id_no           string `json:"Id_no"`
	Barcode_id      string `json:"Barcode_id"`
	Package_date    string `json:"Package_date"`
	Part            string `json:"Part"`
	Weight          string `json:"Weight"`
	Purchase_nm     string `json:"Purchase_nm"`
	Purchase_biz_no string `json:"Purchase_biz_no"`
}

//args[0]						-- Cow Key
//args[1] id_no
//args[2] barcode_id
//args[3] package_date
//args[4] part
//args[5] weight
//args[6] purchase_nm
//args[7] purchase_biz_no

type Remark struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

/*
 * The Init method is called when the Smart Contract "fabcow" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcow"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately

	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAllCows" {
		return s.queryAllCows(APIstub)
	} else if function == "queryAllOwners" {
		return s.queryAllOwners(APIstub)
	} else if function == "registerCow" {
		return s.registerCow(APIstub, args)
	} else if function == "registerHACCP" {
		return s.registerHACCP(APIstub, args)
	} else if function == "registerRFID" {
		return s.registerRFID(APIstub, args)
	} else if function == "registerOwner" {
		return s.registerOwner(APIstub, args)
	} else if function == "registerInProcessesBundleNum" {
		return s.registerInProcessesBundleNum(APIstub, args)
	} else if function == "registerInSalesBundleNum" {
		return s.registerInSalesBundleNum(APIstub, args)
	} else if function == "changeCowOwner" {
		return s.changeCowOwner(APIstub, args)
	} else if function == "query" {
		return s.query(APIstub, args)
	} else if function == "addRemark" {
		return s.addRemark(APIstub, args)
	} else if function == "addBTVaccine" {
		return s.addBTVaccine(APIstub, args)
	} else if function == "addFAMDVaccine" {
		return s.addFAMDVaccine(APIstub, args)
	} else if function == "addInfoDead" {
		return s.addInfoDead(APIstub, args)
	} else if function == "addInfoInspect" {
		return s.addInfoInspect(APIstub, args)
	} else if function == "addInfoGradeResult" {
		return s.addInfoGradeResult(APIstub, args)
	} else if function == "addInfoInProcessesReportPurchase" {
		return s.addInfoInProcessesReportPurchase(APIstub, args)
	} else if function == "addInfoReportPacking" {
		return s.addInfoReportPacking(APIstub, args)
	} else if function == "addInfoReportSale" {
		return s.addInfoReportSale(APIstub, args)
	} else if function == "addInfoInSalesReportPurchase" {
		return s.addInfoInSalesReportPurchase(APIstub, args)
	} else if function == "deleteCow" {
		return s.deleteCow(APIstub, args)
	} else if function == "addAut" {
		return s.addAut(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.(:).)")
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{"Args":["query", "COW", "COW0"]}'
	//args[0]
	//args[1]

	log.Println("--==query==--")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2(:).)")
	}
	if strings.Contains(args[0], "COW") {
		cowAsBytes, _ := APIstub.GetState(args[1])
		return shim.Success(cowAsBytes)
	} else if strings.Contains(args[0], "OWNER") {
		ownerAsBytes, _ := APIstub.GetState(args[1])
		return shim.Success(ownerAsBytes)
	} else if strings.Contains(args[0], "HACCP") {
		haccpAsBytes, _ := APIstub.GetState(args[1])
		return shim.Success(haccpAsBytes)
	} else if strings.Contains(args[0], "RFID") {
		haccpAsBytes, _ := APIstub.GetState(args[1])
		return shim.Success(haccpAsBytes)
	} else if strings.Contains(args[0], "BUNDLE") {
		haccpAsBytes, _ := APIstub.GetState(args[1])
		return shim.Success(haccpAsBytes)
	} else {
		return shim.Error(":)")
	}
}

//���� �� ���� ��������
func (s *SmartContract) queryAllCows(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "COW0"
	endKey := "COW999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- queryAllCows:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//���� ������ ���� ��������
func (s *SmartContract) queryAllOwners(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "OWNER0"
	endKey := "OWNER999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
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

	fmt.Printf("- queryAllCows:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//�⺻ ��, ������ ���� ������(Sample)
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {

	log.Println("--==initLedger==--")

	owners := []Owner{
		Owner{Owner_id: "01", Owner_nm: "ChukLim1", Owner_addr: "Iksan", Livestock: "C", Owner_user_nm: "Kim Duck Bae", Owner_user_birth: "530118"},
		Owner{Owner_id: "02", Owner_nm: "ChukLim2", Owner_addr: "Jeonju", Livestock: "C", Owner_user_nm: "Kim Sam Sun", Owner_user_birth: "520202"},
		Owner{Owner_id: "03", Owner_nm: "ChukLim3", Owner_addr: "Daejeon", Livestock: "C", Owner_user_nm: "Kim Young Mi", Owner_user_birth: "610118"},
	}

	cows := []Cow{
		Cow{Id_no: "180501-2", Birth_date: "180501", Sex: "F", Father_id: "901027", Mother_id: "910101", Origin: "Korea Jeonbuk", Owner: owners[0]},
		Cow{Id_no: "180502-1", Birth_date: "180502", Sex: "M", Father_id: "901027", Mother_id: "910101", Origin: "Korea Jeonbuk", Owner: owners[1]},
		Cow{Id_no: "180503-1", Birth_date: "180503", Sex: "M", Father_id: "901027", Mother_id: "910101", Origin: "Korea Jeonbuk", Owner: owners[2]},
	}

	i := 0
	for i < len(cows) {
		fmt.Println("i is ", i)
		cowAsBytes, _ := json.Marshal(cows)
		APIstub.PutState("COW"+strconv.Itoa(i), cowAsBytes)
		fmt.Println("Added", cows[i])
		i = i + 1
	}

	return shim.Success(nil)
}

//��ü�ĺ���ȣ ����
func (s *SmartContract) registerCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{"Args":["registerCow","COW0", "180501-1","180501", "M", "630118-1", "630331-2", "Ik-San", "OWNER0"]}'
	//args[0]				-- COW Key
	//args[1] Id_no			-- ��ü�ĺ���ȣ
	//args[2] Birth_date	-- ����������
	//args[3] Sex			-- ����
	//args[4] Father_id		-- ݫ ��ü�ĺ���ȣ
	//args[5] Mother_id		-- ٽ ��ü�ĺ���ȣ
	//args[6] Origin		-- ������
	//args[7]				-- OWNER(������)����

	log.Println("--==registerCow==--")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	ownerAsBytes, _ := APIstub.GetState(args[7])
	if ownerAsBytes == nil {
		return shim.Error("Incorrect value. Owner :).")
	}
	owner := Owner{}

	json.Unmarshal(ownerAsBytes, &owner)

	log.Println("Logging: " + owner.Owner_id + "--" + owner.Owner_nm + "--" + owner.Owner_nm + "==" + owner.Owner_addr + "--" + owner.Livestock + "--" + owner.Owner_user_nm + "--" + owner.Owner_user_birth)

	var cow = Cow{Id_no: args[1], Birth_date: args[2], Sex: args[3], Father_id: args[4], Mother_id: args[5], Origin: args[6], Owner: Owner{Owner_id: owner.Owner_id, Owner_nm: owner.Owner_nm, Owner_addr: owner.Owner_addr, Livestock: owner.Livestock, Owner_user_nm: owner.Owner_user_nm, Owner_user_birth: owner.Owner_user_birth}}
	log.Println("Logging: " + cow.Id_no + "--" + cow.Birth_date + "--" + cow.Sex + "==" + cow.Owner.Owner_id + "--" + cow.Owner.Owner_user_nm)

	cowAsBytes, _ := json.Marshal(cow)
	jsonString := string(cowAsBytes)
	log.Println("Logging: " + jsonString)
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//���� ��������, ����������, ������ ����, �Ǹ��� ����
func (s *SmartContract) registerOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["registerOwner","OWNER10", "FARM0", "ChukLim1", "Iksan", "C", "Kim Duck Bae", "530118"]}'
	//'{'"Args":["registerOwner","OWNER11", "SLAUGHTER0", "DoChuk1", "Jeonju", "C", "Lee Do Chuk", "500118", "063-111-2222", "1-7474-8700"]}'
	//'{'"Args":["registerOwner","OWNER12", "PROCESS0", "Gagong1", "PyeongTak", "Empty", "Park Ga Gong", "Empty", "2-7474-8701"]}'
	//'{'"Args":["registerOwner","OWNER13", "SALE0", "Panmae1", "Ansan", "Empty", "Moon Pan Mae", "Empty", "3-7474-8702"]}'
	//args[0]					-- �����ھ��̵�
	///��������(Default)
	//args[1] farm_id			-- �����ĺ���ȣ(slaughter_id[������ID], process_id[������ID], sale_id[�Ǹ���ID])
	//args[2] farm_nm			-- ������(slaughter_nm[��������], process_nm[������ ��ȣ], sale_nm[�Ǹ��� ��ȣ])
	//args[3] farm_addr			-- ����������(slaughter_addr[�������ּ�], process_addr[������ �ּ�], sale_addr[�Ǹ��� �ּ�])
	//args[4] livestock			-- ������ ����(handel_livestock[��������ȭ��ȣ], "Empty", "Empty")
	//args[5] farm_user_nm		-- �����濵�� ����(slaughter_user_nm[������], process_user_nm[������ ��ǥ�ڸ�], sale_user_nm[�Ǹ��� ��ǥ�ڸ�])
	//args[6] farm_user_brith	-- �����濵�� ��������(slaughter_user_birth, "Empty", "Empty")

	///����������(Default ���� ��)
	//args[7] slaughter_tel		-- ��������ȭ��ȣ
	//args[8] slaughter_reg_no	-- �����������ڵ��Ϲ�ȣ

	///����������(Default ���� ��)
	//args[7] process_biz_no	-- ������ �����ڹ�ȣ

	///�Ǹ�������(Default ���� ��)
	//args[7] sale_biz_no		-- �Ǹ��� �����ڹ�ȣ

	log.Println("--==registerOwner==--")

	//�ĺ���ȣ Ȯ��
	if strings.Contains(args[1], "FARM") {
		log.Println("--==>>registerOwner[FARM]")
		//�Ķ����� Ȯ��
		if len(args) != 7 {
			return shim.Error("Incorrect number of arguments. Expecting 7")
		}

		//struct ������ ����
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		log.Println("Logging: " + owner.Owner_id + "--" + owner.Owner_nm + "--" + owner.Livestock + "==" + owner.Owner_user_nm + "--" + owner.Owner_user_birth)

		//JSON Ÿ������ ����
		ownerAsBytes, _ := json.Marshal(owner)
		jsonString := string(ownerAsBytes)
		log.Println("Logging: " + jsonString)
		APIstub.PutState(args[0], ownerAsBytes)

		return shim.Success(nil)
	} else if strings.Contains(args[1], "SLAUGHTER") {
		log.Println("--==>>registerOwner[SLAUGHTER]")
		if len(args) != 9 {
			return shim.Error("Incorrect number of arguments. Expecting 9")
		}

		//struct�� ����
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		//Default ������ �� �ܿ� �ٸ� ���� ����
		variables := [2]string{"registerOwner.slaughter_tel", "registerOwner.slaughter_reg_no"}
		log.Println(variables)
		for i := 0; i < len(variables); i++ {
			log.Println("For Loop")
			log.Println(i)
			remarkData := Remark{Key: variables[i], Value: args[i+2]}
			log.Println(remarkData)
			owner.setOwnerRemark(remarkData)
		}

		log.Println("Logging: " + owner.Owner_id + " -- " + owner.Owner_nm + " -- " + owner.Livestock + " -- " + owner.Owner_user_nm + " -- " + owner.Owner_user_birth + " -- ")

		//JSON Ÿ������ ����
		ownerAsBytes, _ := json.Marshal(owner)
		jsonString := string(ownerAsBytes)
		log.Println("Logging: " + jsonString)
		APIstub.PutState(args[0], ownerAsBytes)

		return shim.Success(nil)

	} else if strings.Contains(args[1], "PROCESS") {
		log.Println("--==>>registerOwner[PROCESS]")

		if len(args) != 8 {
			return shim.Error("Incorrect number of arguments. Expecting 8")
		}

		//struct�� ����
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		//Default ������ �� �ܿ� �ٸ� ���� ����
		variables := [1]string{"registerOwner.process_biz_no"}
		log.Println(variables)
		for i := 0; i < len(variables); i++ {
			log.Println("For Loop")
			log.Println(i)
			remarkData := Remark{Key: variables[i], Value: args[i+2]}
			log.Println(remarkData)
			owner.setOwnerRemark(remarkData)
		}

		log.Println("Logging: " + owner.Owner_id + " -- " + owner.Owner_nm + " -- " + owner.Livestock + " -- " + owner.Owner_user_nm + " -- " + owner.Owner_user_birth + " -- ")

		//JSON Ÿ������ ����
		ownerAsBytes, _ := json.Marshal(owner)
		jsonString := string(ownerAsBytes)
		log.Println("Logging: " + jsonString)
		APIstub.PutState(args[0], ownerAsBytes)

		return shim.Success(nil)

	} else if strings.Contains(args[1], "SALE") {
		log.Println("--==>>registerOwner[SALE]")

		if len(args) != 8 {
			return shim.Error("Incorrect number of arguments. Expecting 8")
		}

		//struct�� ����
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		//Default ������ �� �ܿ� �ٸ� ���� ����
		variables := [1]string{"registerOwner.sale_biz_no"}
		log.Println(variables)
		for i := 0; i < len(variables); i++ {
			log.Println("For Loop")
			log.Println(i)
			remarkData := Remark{Key: variables[i], Value: args[i+2]}
			log.Println(remarkData)
			owner.setOwnerRemark(remarkData)
		}

		log.Println("Logging: " + owner.Owner_id + " -- " + owner.Owner_nm + " -- " + owner.Livestock + " -- " + owner.Owner_user_nm + " -- " + owner.Owner_user_birth + " -- ")

		//JSON Ÿ������ ����
		ownerAsBytes, _ := json.Marshal(owner)
		jsonString := string(ownerAsBytes)
		log.Println("Logging: " + jsonString)
		APIstub.PutState(args[0], ownerAsBytes)

		return shim.Success(nil)

	}

	return shim.Success(nil)
}

//HACCP ��������
func (s *SmartContract) registerHACCP(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["registerHACCP","HACCP0", "OWNER10","FARM0", "ChukLim1", "Iksan", "Cow", "20280528"]}'
	//HACCP ����, HACCP���̵�, �����ھ��̵�, ...
	//args[0]				-- HACCP���̵�
	//args[1]				-- �����ھ��̵�
	//args[2] farm_id		-- �������̵�
	//args[3] farm_nm		-- ������
	//args[4] farm_addr		-- �����ּ�
	//args[5] apply_item	-- ����ǰ��
	//args[6] validity_date	-- ��ȿ�Ⱓ

	log.Println("--==registerHACCP==--")

	//���� Ȯ��
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//HACCP INVOKE
	var haccp = HACCP{Farm_id: args[2], Farm_nm: args[3], Farm_addr: args[4], Apply_item: args[5], Validity_date: args[6]}
	log.Println("Logging: " + haccp.Farm_id + "--" + haccp.Farm_nm + "--" + haccp.Farm_addr + "==" + haccp.Apply_item + "--" + haccp.Validity_date)

	haccpAsBytes, _ := json.Marshal(haccp)
	jsonString := string(haccpAsBytes)
	log.Println("Logging: " + jsonString)
	//HACCP Asset ����
	APIstub.PutState(args[0], haccpAsBytes)

	//OWNER INVOKE
	ownerAsBytes, _ := APIstub.GetState(args[1])
	owner := Owner{}
	json.Unmarshal(ownerAsBytes, &owner)

	variables := [2]string{"haccp.Id_no", "haccp.Rfid_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+2]}
		log.Println(remarkData)
		owner.setOwnerRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	ownerAsBytes, _ = json.Marshal(owner)
	//PutState����
	APIstub.PutState(args[1], ownerAsBytes)

	return shim.Success(nil)
}

//��ǥ ���� - �ڻ�
func (s *SmartContract) registerRFID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["registerRFID", "COW3", "RFID0"]}'
	//HACCP ����, HACCP���̵�, �����ھ��̵�, ...
	//args[0]				-- ��ü �ĺ���ȣ
	//args[1]				-- RFID �ĺ���ȣ

	log.Println("--==registerRFID==--")

	//���� Ȯ��
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//RFID INVOKE
	var rfid = RFID{Id_no: args[0], Rfid_no: args[1]}
	log.Println("Logging: " + rfid.Id_no + "--" + rfid.Rfid_no)

	rfidAsBytes, _ := json.Marshal(rfid)
	jsonString := string(rfidAsBytes)
	log.Println("Logging: " + jsonString)

	//RFID Asset ����
	APIstub.PutState(args[1], rfidAsBytes)

	//COW INVOKE
	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)

	variables := [2]string{"rfid.Id_no", "rfid.Rfid_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//������ȣ ���� - �ڻ�(������)
func (s *SmartContract) registerInProcessesBundleNum(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//�����幭����ȣ����

	log.Println("--==registerInProcessesBundleNum==--")

	//'{'"Args":["registerInProcessesBundleNum","BUNDLE0", "COW10", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- ��ü�ĺ���ȣ
	//args[2] barcode_id			-- ���ڵ� ID
	//args[3] package_date			-- ����ó����
	//args[4] part					-- ����
	//args[5] weight				-- �߷�
	//args[6] purchase_nm			-- ����/�Ƿ�ó ��ȣ
	//args[7] purchase_biz_no		-- ����ó �����ڵ��Ϲ�ȣ

	//���� Ȯ��
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//Bundle INVOKE
	var bundle = Bundle{Id_no: args[1], Barcode_id: args[2], Package_date: args[3], Part: args[4], Weight: args[5], Purchase_nm: args[6], Purchase_biz_no: args[7]}

	bundleAsBytes, _ := json.Marshal(bundle)
	jsonString := string(bundleAsBytes)
	log.Println("Logging: " + jsonString)
	//Bundle Asset ����
	APIstub.PutState(args[0], bundleAsBytes)

	//COW INVOKE
	cowAsBytes, _ := APIstub.GetState(args[1])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)

	variables := [7]string{"registerInProcessesBundleNum.id_no", "registerInProcessesBundleNum.barcode_id", "registerInProcessesBundleNum.package_date", "registerInProcessesBundleNum.part", "registerInProcessesBundleNum.weight", "registerInProcessesBundleNum.purchase_nm", "registerInProcessesBundleNum.purchase_biz_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+2]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[1], cowAsBytes)

	return shim.Success(nil)
}

//������ȣ ���� - �ڻ�(�Ǹ���)
func (s *SmartContract) registerInSalesBundleNum(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//�Ǹ��幭����ȣ����

	log.Println("--==registerInSalesBundleNum==--")

	//'{'"Args":["registerInSalesBundleNum","BUNDLE0", "COW10", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- ��ü�ĺ���ȣ
	//args[2] barcode_id			-- ���ڵ� ID
	//args[3] package_date			-- ����ó����
	//args[4] part					-- ����
	//args[5] weight				-- �߷�
	//args[6] purchase_nm			-- ����/�Ƿ�ó ��ȣ
	//args[7] purchase_biz_no		-- ����ó �����ڵ��Ϲ�ȣ

	//���� Ȯ��
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//Bundle INVOKE
	var bundle = Bundle{Id_no: args[1], Barcode_id: args[2], Package_date: args[3], Part: args[4], Weight: args[5], Purchase_nm: args[6], Purchase_biz_no: args[7]}

	bundleAsBytes, _ := json.Marshal(bundle)
	jsonString := string(bundleAsBytes)
	log.Println("Logging: " + jsonString)
	//Bundle Asset ����
	APIstub.PutState(args[0], bundleAsBytes)

	//COW INVOKE
	cowAsBytes, _ := APIstub.GetState(args[1])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)

	variables := [7]string{"registerInSalesBundleNumbundle.id_no", "registerInSalesBundleNum.barcode_id", "registerInSalesBundleNum.package_date", "registerInSalesBundleNum.part", "registerInSalesBundleNum.weight", "registerInSalesBundleNum.purchase_nm", "registerInSalesBundleNum.purchase_biz_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+2]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[1], cowAsBytes)

	return shim.Success(nil)
}

//���� �̵�(����/�絵), ������ �̵�, ������ �̵�, �Ǹ��� �̵�
func (s *SmartContract) changeCowOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{"Args":["changeCowOwner", "COW3", "OWNER1", "OWNER2"]}'
	//args[0]				-- �Ҿ��̵�
	//args[1]				-- �����Ҽ����ھ��̵�
	//args[2]				-- �絵�ɼ����ھ��̵�

	log.Println("--==changeCowOwner==--")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	//OWNER 1,2,3�� ����
	json.Unmarshal(cowAsBytes, &cow)

	// if cow.Owner.Owner_id != args[1] {
	// 	return shim.Error("�絵�� ���� ������ ���̵��� �ùٸ��� �ʽ��ϴ�.")
	// }

	//�絵���� ������ ���� ��������
	ownerAsBytes, _ := APIstub.GetState(args[2])
	owner := Owner{}

	json.Unmarshal(ownerAsBytes, &owner)

	cow.Owner.Owner_id = owner.Owner_id
	cow.Owner.Owner_nm = owner.Owner_nm
	cow.Owner.Owner_addr = owner.Owner_addr
	cow.Owner.Livestock = owner.Livestock
	cow.Owner.Owner_user_nm = owner.Owner_user_nm
	cow.Owner.Owner_user_birth = owner.Owner_user_birth
	cow.Owner.Remarks = owner.Remarks

	cowAsBytes, _ = json.Marshal(cow)
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//ģȯ�� ��ǰ��������
func (s *SmartContract) addAut(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addAut==--")

	//'{'"Args":["addAut","OWNER90", "AutCheck", "201905251610", "�า", "1985.10.27", "�����ϵ� �ͻ���", "Jeon Buk", "Cow", "30", "������", "10001", "20180525"]}'
	//aut_falg 			-- ��������
	//validity_date		-- ��ȿ�Ⱓ
	//farm_nm			-- ��������
	//farm_birth_date	-- �����ڻ�������
	//farm_addr			-- ������ �ּ�
	//biz_addr			-- ������ ������
	//aut_item			-- ����ǰ��
	//breed_head		-- �����μ�
	//aut_com			-- ��������
	//aut_id			-- ������ȣ
	//aut_date			-- ��������

	// if len(args) != 13 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 13")
	// }

	log.Println(len(args))

	ownerAsBytes, _ := APIstub.GetState(args[0])
	owner := Owner{}
	json.Unmarshal(ownerAsBytes, &owner)
	log.Println("check Owner bottom line ___")
	log.Println(owner)

	variables := [11]string{"aut_falg", "validity_date", "farm_nm", "farm_birth_date", "farm_addr", "biz_addr", "aut_item", "breed_head", "aut_com", "aut_id", "aut_date"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		owner.setOwnerRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	ownerAsBytes, _ = json.Marshal(owner)
	//PutState����
	APIstub.PutState(args[0], ownerAsBytes)

	return shim.Success(nil)
}

//����/���缿��(Tuberculousis/Brucella) �˻� - �ŷ�
func (s *SmartContract) addBTVaccine(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addBTVaccine==--")

	//'{'"Args":["addBTVaccine","COW90", "farm_id", "farm_nm", "farm_addr", "farm_user_nm", "farm_user_birth", "farm_user_addr", "inspection_date", "inspection_head", "inspection_method", "livestock", "kind", "sex", "age", "id_no", "inspection_result", "inspection_part", "inspection_user_nm"]}'
	//args[0]						-- �Ҿ��̵�
	//args[1] 	farm_id				-- �����ĺ���ȣ
	//args[2] 	farm_nm				-- ������
	//args[3] 	farm_addr			-- ������
	//args[4] 	farm_user_nm		-- ���������� ����
	//args[5] 	farm_user_birth		-- ���������� ��������
	//args[6] 	farm_user_addr		-- ���������� �ּ�
	//args[7] 	inspection_date		-- �˻翬����
	//args[8] 	inspection_head		-- �˻��μ�
	//args[9] 	inspection_method	-- �˻�����
	//args[10]	livestock			-- ����
	//args[11] 	kind				-- ǰ��
	//args[12] 	sex					-- ����
	//args[13] 	age					-- ����
	//args[14] 	id_no				-- ��ü�ĺ���ȣ
	//args[15] 	inspection_result	-- �˻�����
	//args[16] 	inspection_part		-- �˻��ڼҼ�
	//args[17] 	inspection_user_nm	-- �˻��ڼ���

	if len(args) != 18 {
		return shim.Error("Incorrect number of arguments. Expecting 18")
	}

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Cow bottom line ___")
	log.Println(cow)

	variables := [17]string{"addBTVaccine.farm_id", "addBTVaccine.farm_nm", "addBTVaccine.farm_addr", "addBTVaccine.farm_user_nm", "addBTVaccine.farm_user_birth", "addBTVaccine.farm_user_addr", "addBTVaccine.inspection_date", "addBTVaccine.inspection_head", "addBTVaccine.inspection_method",
		"addBTVaccine.livestock", "addBTVaccine.kind", "addBTVaccine.sex", "addBTVaccine.age", "addBTVaccine.id_no", "addBTVaccine.inspection_result", "addBTVaccine.inspection_part", "addBTVaccine.inspection_user_nm"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//������ ��������(Foot And Mouse Disease)
func (s *SmartContract) addFAMDVaccine(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addFAMDVaccine==--")

	//'{'"Args":["addFAMDVaccine","COW90", "farm_id", "farm_addr", "farm_tel", "breed_head", "item", "sex", "age", "id_no", "vaccination_date"]}'
	//args[0]					-- �Ҿ��̵�
	//args[1] farm_id			-- �����ĺ���ȣ
	//args[2] farm_addr			-- �����ּ�
	//args[3] farm_tel			-- ������ȭ��ȣ
	//args[4] breed_head		-- �����μ�
	//args[5] item				-- ǰ��
	//args[6] sex				-- ����
	//args[7] age				-- ����
	//args[8] id_no				-- ��ü�ĺ���ȣ
	//args[9] vaccination_date	-- ����������

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Cow bottom line ___")
	log.Println(cow)

	variables := [9]string{"addFAMDVaccine.farm_id", "addFAMDVaccine.farm_addr", "addFAMDVaccine.farm_tel", "addFAMDVaccine.breed_head", "addFAMDVaccine.item", "addFAMDVaccine.sex", "addFAMDVaccine.age", "addFAMDVaccine.id_no", "addFAMDVaccine.vaccination_date"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//���� ���� - �ŷ�
func (s *SmartContract) addInfoDead(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoDead==--")

	//'{'"Args":["addInfoDead","COW90", "01", "180501-1", "180528", "Cancer", "burning"]}'
	//args[0]				-- �Ҿ��̵�
	//args[1] farm_id		-- �������̵�
	//args[2] id_no			-- ��ü�ĺ���ȣ
	//args[3] det_date		-- ������
	//args[4] det_reason	-- ��������
	//args[5] det_method	-- ����ó������

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [5]string{"addInfoDead.farm_id", "addInfoDead.id_no", "addInfoDead.det_date", "addInfoDead.det_reason", "addInfoDead.det_method"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//�������� - �ŷ�
func (s *SmartContract) addInfoDeliver(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoDeliver==--")

	//'{'"Args":["addInfoDeliver","COW10", "180501-1", "RFID0"]}'
	//args[0]				-- Cow Key
	//args[1] id_no			-- ��ü�ĺ���ȣ
	//args[2] rfid_no		-- RFID�ĺ���ȣ

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [2]string{"addInfoDeliver.id_no", "addInfoDeliver.rfid_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//�����˻����� - �ŷ�
func (s *SmartContract) addInfoInspect(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoInspect==--")

	//'{'"Args":["addInfoInspect","COW10", "COW", "180501-1", "300kg", "DoChuk1", "seal_10", "20180529", "FARM0", "Jeonju", "HACCP10", "Discard", "20180529", "Korea Inspect Center", "vetrinarian_100"]}'
	//args[0]						-- Cow Key
	//args[1] livestock				-- ������ ����
	//args[2] id_no					-- ��ü�ĺ���ȣ
	//args[3] weight				-- �߷�(���� �� ����)
	//args[4] slaughter_nm			-- ��������
	//args[5] seal_no				-- ���ι�ȣ
	//args[6] slaughter_date		-- ���࿬����
	//args[7] farm_id				-- �����ĺ���ȣ
	//args[8] farm_addr				-- �����Ƿ��� �ּ�
	//args[9] haccp_yn				-- HACCP ��������
	//args[10] fale_method			-- ���հ� ���� ó�� ����
	//args[11] inspection_date		-- �˻翬����
	//args[12] inspection_part		-- �˻��� �Ҽ�
	//args[13] inspection_user_nm	-- �˻��� ����
	//args[14] veterinarian_no		-- ���ǻ� ������ȣ

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [14]string{"addInfoInspect.livestock", "addInfoInspect.id_no", "addInfoInspect.weight", "addInfoInspect.slaughter_nm", "addInfoInspect.seal_no", "addInfoInspect.slaughter_date", "addInfoInspect.farm_id", "addInfoInspect.farm_addr", "addInfoInspect.haccp_yn", "addInfoInspect.fale_method", "addInfoInspect.inspection_date", "addInfoInspect.inspection_part", "addInfoInspect.inspection_user_nm", "addInfoInspect.veterinarian_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//���������������� - �ŷ�
func (s *SmartContract) addInfoGradeResult(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoGradeResult==--")

	//'{'"Args":["addInfoGradeResult","COW10", "180529", "quality_part", "quality_nm", "subscriber_nm", "subscriber_birth", "subscriber_company", "subscriber_addr", "slaughter_nm", "slaughter_addr", "id_no", "weight", "meat_quality_grade", "meat_weight_grade", "grade_head"]}'
	//args[0]						-- Cow Key
	//args[1] grade_date			-- ��������������
	//args[2] quality_part			-- ǰ���򰡻� �Ҽ�
	//args[3] quality_nm			-- ǰ���򰡻� ����
	//args[4] subscriber_nm			-- ��û�� ����
	//args[5] subscriber_birth		-- ��û�� ��������
	//args[6] subscriber_company	-- ��û�� ���Ҹ�
	//args[7] subscriber_addr		-- ��û�� �ּ�
	//args[8] slaughter_nm			-- ������ ��
	//args[9] slaughter_addr		-- ������ �ּ�
	//args[10] id_no				-- ��ü�ĺ���ȣ
	//args[11] weight				-- �߷�(���� �� ����)
	//args[12] meat_quality_grade	-- ��������
	//args[13] meat_weight_grade	-- ��������
	//args[14] grade_head			-- �����μ�

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [14]string{"addInfoGradeResult.grade_date", "addInfoGradeResult.quality_part", "addInfoGradeResult.quality_nm", "addInfoGradeResult.subscriber_nm", "addInfoGradeResult.subscriber_birth", "addInfoGradeResult.subscriber_company", "addInfoGradeResult.subscriber_addr", "addInfoGradeResult.slaughter_nm", "addInfoGradeResult.slaughter_addr", "addInfoGradeResult.id_no", "addInfoGradeResult.weight", "addInfoGradeResult.meat_quality_grade", "addInfoGradeResult.meat_weight_grade", "addInfoGradeResult.grade_head"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//���ԽŰ� - �ŷ�(������)
func (s *SmartContract) addInfoInProcessesReportPurchase(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//���������ԽŰ�

	log.Println("--==addInfoInProcessesReportPurchase==--")

	//'{'"Args":["addInfoInProcessesReportPurchase","COW10", "barcode_id", "deal_date", "origin", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] barcode_id			-- ���ڵ�ID
	//args[2] deal_date				-- �ŷ�������
	//args[3] origin				-- ������
	//args[4] part					-- ����
	//args[5] weight				-- �߷�
	//args[6] purchase_nm			-- ����ó ��ȣ
	//args[7] purchase_biz_no		-- ����ó �����ڵ��Ϲ�ȣ

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [7]string{"addInfoInProcessesReportPurchase.barcode_id", "addInfoInProcessesReportPurchase.deal_date", "addInfoInProcessesReportPurchase.origin", "addInfoInProcessesReportPurchase.part", "addInfoInProcessesReportPurchase.weight", "addInfoInProcessesReportPurchase.purchase_nm", "addInfoInProcessesReportPurchase.purchase_biz_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)
	return shim.Success(nil)
}

//����ó�������Ű� - �ŷ�
func (s *SmartContract) addInfoReportPacking(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoReportPacking==--")

	//'{'"Args":["addInfoReportPacking","COW10", "id_no", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- ��ü�ĺ���ȣ
	//args[2] barcode_id			-- ���ڵ� ID
	//args[3] package_date			-- ����ó����
	//args[4] part					-- ����
	//args[5] weight				-- �߷�
	//args[6] purchase_nm			-- ����/�Ƿ�ó ��ȣ
	//args[7] purchase_biz_no		-- ����ó �����ڵ��Ϲ�ȣ

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [7]string{"addInfoReportPacking.id_no", "addInfoReportPacking.barcode_id", "addInfoReportPacking.package_date", "addInfoReportPacking.part", "addInfoReportPacking.weight", "addInfoReportPacking.purchase_nm", "addInfoReportPacking.purchase_biz_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)
	return shim.Success(nil)
}

//�ǸŽŰ� - �ŷ�
func (s *SmartContract) addInfoReportSale(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoReportSale==--")

	//'{'"Args":["addInfoReportSale","COW10", "id_no", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- ��ü�ĺ���ȣ
	//args[2] barcode_id			-- ���ڵ� ID
	//args[3] sale_date				-- �Ǹſ�����
	//args[4] part					-- ����
	//args[5] weight				-- �Ǹ��߷�
	//args[6] sale_nm				-- �Ǹ�ó ��ȣ
	//args[7] sale_biz_no			-- �Ǹ�ó �����ڵ��Ϲ�ȣ

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [7]string{"addInfoReportSale.id_no", "addInfoReportSale.barcode_id", "addInfoReportSale.sale_date", "addInfoReportSale.part", "addInfoReportSale.weight", "addInfoReportSale.sale_nm", "addInfoReportSale.sale_biz_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//���ԽŰ� - �ŷ�(�Ǹ���)
func (s *SmartContract) addInfoInSalesReportPurchase(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoInSalesReportPurchase==--")

	//'{'"Args":["addInfoInSalesReportPurchase","COW10", "barcode_id", "deal_date", "origin", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] barcode_id			-- ���ڵ�ID
	//args[2] deal_date				-- �ŷ�������
	//args[3] origin				-- ������
	//args[4] part					-- ����
	//args[5] weight				-- �߷�
	//args[6] purchase_nm			-- ����ó ��ȣ
	//args[7] purchase_biz_no		-- ����ó �����ڵ��Ϲ�ȣ

	log.Println(len(args))

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	json.Unmarshal(cowAsBytes, &cow)
	log.Println("check Owner bottom line ___")
	log.Println(cow)

	variables := [7]string{"addInfoInSalesReportPurchase.barcode_id", "addInfoInSalesReportPurchase.deal_date", "addInfoInSalesReportPurchase.origin", "addInfoInSalesReportPurchase.part", "addInfoInSalesReportPurchase.weight", "addInfoInSalesReportPurchase.purchase_nm", "addInfoInSalesReportPurchase.purchase_biz_no"}
	log.Println(variables)
	for i := 0; i < len(variables); i++ {
		log.Println("For Loop")
		log.Println(i)
		remarkData := Remark{Key: variables[i], Value: args[i+1]}
		log.Println(remarkData)
		cow.setCowRemark(remarkData)
	}

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)
	return shim.Success(nil)
}

//�� ���� ����
func (s *SmartContract) deleteCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string
	var CowJSON Cow
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	cowId := args[0]

	// to maintain the make~grade index, we need to read the cow first and get 'Make' values.
	// make ~ grade �ε����� �����Ϸ��� ���� ���� ������ �а� '������ü' ������ �����;��մϴ�.
	valAsbytes, err := APIstub.GetState(cowId) //get the cow from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cowId + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Cow does not exist: " + cowId + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &CowJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + cowId + "\"}"
		return shim.Error(jsonResp)
	}

	err = APIstub.DelState(cowId) //remove the cow from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// // maintain the index
	// indexName := "color~name"
	// colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{CowJSON.Make, CowJSON.Model})
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }

	// //  Delete index entry to state.
	// err = APIstub.DelState(colorNameIndexKey)
	// if err != nil {
	// 	return shim.Error("Failed to delete state:" + err.Error())
	// }
	return shim.Success(nil)
}

//test(Remark ���� Test)
func (s *SmartContract) addRemark(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["addRemark","COW4", "addVaccine", "True"]}

	log.Println("--==addRemark==--")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	log.Println("args[0]: " + args[0] + "  args[1]: " + args[1] + "  args[2]: " + args[2])
	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}

	//Byte������ json�� ���� json���� ����
	json.Unmarshal(cowAsBytes, &cow)

	//�� ���������� ���� Key Value�� ����
	remark1 := Remark{Key: args[1], Value: args[2]}

	//�������� ���� ������ ����
	//Remarks := []Remark{}

	cow.setCowRemark(remark1)

	// if cow.Remarks == nil {
	// 	cow.Remarks[0] = remark1
	// } else {
	// 	//������ Remarks �迭�� �� �������� �߰�
	// 	Remarks = append(cow.Remarks, remark1)
	// }
	log.Println("change [cow] cow.Remarks Key: " + cow.Remarks[0].Key + " Value: " + cow.Remarks[0].Value)

	//������ Remark ���� �߰��� Remarks�� �ٽ� �־���
	//cow.Remarks = Remarks

	//Json�� �ٽ� Byte ���·� ����
	cowAsBytes, _ = json.Marshal(cow)
	//PutState����
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//test(Cow Remark ���� Test)
func (cow *Cow) setCowRemark(remark1 Remark) []Remark {
	cow.Remarks = append(cow.Remarks, remark1)
	return cow.Remarks
}

//test(Owner Remark ���� Test)
func (owner *Owner) setOwnerRemark(remark1 Remark) []Remark {
	owner.Remarks = append(owner.Remarks, remark1)
	return owner.Remarks
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
