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
 * 설명서 항목 샘플 스마트 계약 :
 * 첫 블록 체인 어플리케이션 작성하기
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 * 포맷팅, 바이트 처리, JSON 읽기 및 쓰기, 문자열 조작을 위한 유틸리티 라이브러리 4
 * 스마트 계약을 위한 특정 Hyperbelger Fabric 라이브러리 2 개
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
// 스마트 컨트렉 구조를 정의
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
// 7가지 특성들을 가진 소의 구조를 정의합니다. 구조 태그들은 인코딩 혹은 json 라이브러리에 의해 사용됩니다.
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
//args[1] id_no					-- 개체식별번호
//args[2] barcode_id			-- 바코드 ID
//args[3] package_date			-- 포장처리일
//args[4] part					-- 부위
//args[5] weight				-- 중량
//args[6] purchase_nm			-- 매입/의뢰처 상호
//args[7] purchase_biz_no		-- 매입처 사업자등록번호

type Remark struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

/*
 * The Init method is called when the Smart Contract "fabcow" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 * Init 메소드는 블록체인 네트워크에 의해 "fabcow"의 스마트 컨트렉이 인스턴스화될 때 호출됩니다.
 * 가장 좋은 실습은 다른 기능안에서 어떠한 원장이든 초기화를 하는것입니다다. -- initLedger() 함수를 보십시오.
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcow"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 * Invoke 메서드는 스마트 계약 "fabcow"를 실행하기위한 응용 프로그램 요청의 결과로 호출됩니다.
 * 호출 응용 프로그램은 인수를 사용하여 특정 스마트 계약 함수를 호출하도록 지정했습니다.
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	// 요청된 스마트 컨트렉 함수와 인수들을 검색합니다.
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	// 원장과 상호작용하는 처리 함수로 연결합니다.
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

	return shim.Error("Invalid Smart Contract function name.(유효하지 않은 스마트 컨트렉 함수 이름입니다.)")
}

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{"Args":["query", "COW", "COW0"]}'
	//args[0]				-- 자산정보(COW, OWNER, HACCP, RFID)
	//args[1]				-- (자산+넘버링)

	log.Println("--==query==--")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2(적절하지 않은 인수들의 수입니다.)")
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
		return shim.Error("적절하지 않은 인자값 입니다.")
	}
}

//모든 소 정보 가져오기
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

//모든 소유자 정보 가져오기
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

//기본 소, 소유자 정보 만들기(Sample)
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

//개체식별번호 등록
func (s *SmartContract) registerCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{"Args":["registerCow","COW0", "180501-1","180501", "M", "630118-1", "630331-2", "Ik-San", "OWNER0"]}'
	//args[0]				-- COW Key
	//args[1] Id_no			-- 개체식별번호
	//args[2] Birth_date	-- 출생년월일
	//args[3] Sex			-- 성별
	//args[4] Father_id		-- 父 개체식별번호
	//args[5] Mother_id		-- 母 개체식별번호
	//args[6] Origin		-- 원산지
	//args[7]				-- OWNER(농장주)정보

	log.Println("--==registerCow==--")

	if len(args) != 8 {
		return shim.Error("Incorrect number of arguments. Expecting 8")
	}

	ownerAsBytes, _ := APIstub.GetState(args[7])
	if ownerAsBytes == nil {
		return shim.Error("Incorrect value. Owner 정보를 올바르게 입력하셔야 합니다.")
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

//농장 정보등록, 도축장등록, 가공장 등록, 판매장 등록
func (s *SmartContract) registerOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["registerOwner","OWNER10", "FARM0", "ChukLim1", "Iksan", "C", "Kim Duck Bae", "530118"]}'
	//'{'"Args":["registerOwner","OWNER11", "SLAUGHTER0", "DoChuk1", "Jeonju", "C", "Lee Do Chuk", "500118", "063-111-2222", "1-7474-8700"]}'
	//'{'"Args":["registerOwner","OWNER12", "PROCESS0", "Gagong1", "PyeongTak", "Empty", "Park Ga Gong", "Empty", "2-7474-8701"]}'
	//'{'"Args":["registerOwner","OWNER13", "SALE0", "Panmae1", "Ansan", "Empty", "Moon Pan Mae", "Empty", "3-7474-8702"]}'
	//args[0]					-- 소유자아이디
	///농장정보(Default)
	//args[1] farm_id			-- 농장식별번호(slaughter_id[도축장ID], process_id[가공장ID], sale_id[판매장ID])
	//args[2] farm_nm			-- 농장명(slaughter_nm[도축장명], process_nm[가공장 상호], sale_nm[판매장 상호])
	//args[3] farm_addr			-- 농장소재지(slaughter_addr[도축장주소], process_addr[가공장 주소], sale_addr[판매장 주소])
	//args[4] livestock			-- 가축의 종류(handel_livestock[도축장전화번호], "Empty", "Empty")
	//args[5] farm_user_nm		-- 농장경영자 성명(slaughter_user_nm[도축장], process_user_nm[가공장 대표자명], sale_user_nm[판매장 대표자명])
	//args[6] farm_user_brith	-- 농장경영자 생년월일(slaughter_user_birth, "Empty", "Empty")

	///도축장정보(Default 정보 외)
	//args[7] slaughter_tel		-- 도축장전화번호
	//args[8] slaughter_reg_no	-- 도축장사업자등록번호

	///가공장정보(Default 정보 외)
	//args[7] process_biz_no	-- 가공장 사업자번호

	///판매장정보(Default 정보 외)
	//args[7] sale_biz_no		-- 판매장 사업자번호

	log.Println("--==registerOwner==--")

	//식별번호 확인
	if strings.Contains(args[1], "FARM") {
		log.Println("--==>>registerOwner[FARM]")
		//파라미터 확인
		if len(args) != 7 {
			return shim.Error("Incorrect number of arguments. Expecting 7")
		}

		//struct 변수에 대입
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		log.Println("Logging: " + owner.Owner_id + "--" + owner.Owner_nm + "--" + owner.Livestock + "==" + owner.Owner_user_nm + "--" + owner.Owner_user_birth)

		//JSON 타입으로 변경
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

		//struct에 대입
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		//Default 지정된 값 외에 다른 변수 대입
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

		//JSON 타입으로 변경
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

		//struct에 대입
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		//Default 지정된 값 외에 다른 변수 대입
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

		//JSON 타입으로 변경
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

		//struct에 대입
		var owner = Owner{Owner_id: args[1], Owner_nm: args[2], Owner_addr: args[3], Livestock: args[4], Owner_user_nm: args[5], Owner_user_birth: args[6]}

		//Default 지정된 값 외에 다른 변수 대입
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

		//JSON 타입으로 변경
		ownerAsBytes, _ := json.Marshal(owner)
		jsonString := string(ownerAsBytes)
		log.Println("Logging: " + jsonString)
		APIstub.PutState(args[0], ownerAsBytes)

		return shim.Success(nil)

	}

	return shim.Success(nil)
}

//HACCP 인증등록
func (s *SmartContract) registerHACCP(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["registerHACCP","HACCP0", "OWNER10","FARM0", "ChukLim1", "Iksan", "Cow", "20280528"]}'
	//HACCP 등록, HACCP아이디, 소유자아이디, ...
	//args[0]				-- HACCP아이디
	//args[1]				-- 소유자아이디
	//args[2] farm_id		-- 농장아이디
	//args[3] farm_nm		-- 농장명
	//args[4] farm_addr		-- 농장주소
	//args[5] apply_item	-- 적용품목
	//args[6] validity_date	-- 유효기간

	log.Println("--==registerHACCP==--")

	//인자 확인
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//HACCP INVOKE
	var haccp = HACCP{Farm_id: args[2], Farm_nm: args[3], Farm_addr: args[4], Apply_item: args[5], Validity_date: args[6]}
	log.Println("Logging: " + haccp.Farm_id + "--" + haccp.Farm_nm + "--" + haccp.Farm_addr + "==" + haccp.Apply_item + "--" + haccp.Validity_date)

	haccpAsBytes, _ := json.Marshal(haccp)
	jsonString := string(haccpAsBytes)
	log.Println("Logging: " + jsonString)
	//HACCP Asset 등록
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

	//Json을 다시 Byte 형태로 변경
	ownerAsBytes, _ = json.Marshal(owner)
	//PutState해줌
	APIstub.PutState(args[1], ownerAsBytes)

	return shim.Success(nil)
}

//귀표 부착 - 자산
func (s *SmartContract) registerRFID(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["registerRFID", "COW3", "RFID0"]}'
	//HACCP 등록, HACCP아이디, 소유자아이디, ...
	//args[0]				-- 개체 식별번호
	//args[1]				-- RFID 식별번호

	log.Println("--==registerRFID==--")

	//인자 확인
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	//RFID INVOKE
	var rfid = RFID{Id_no: args[0], Rfid_no: args[1]}
	log.Println("Logging: " + rfid.Id_no + "--" + rfid.Rfid_no)

	rfidAsBytes, _ := json.Marshal(rfid)
	jsonString := string(rfidAsBytes)
	log.Println("Logging: " + jsonString)

	//RFID Asset 등록
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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//묶음번호 등록 - 자산(가공장)
func (s *SmartContract) registerInProcessesBundleNum(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//가공장묶음번호등록

	log.Println("--==registerInProcessesBundleNum==--")

	//'{'"Args":["registerInProcessesBundleNum","BUNDLE0", "COW10", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- 개체식별번호
	//args[2] barcode_id			-- 바코드 ID
	//args[3] package_date			-- 포장처리일
	//args[4] part					-- 부위
	//args[5] weight				-- 중량
	//args[6] purchase_nm			-- 매입/의뢰처 상호
	//args[7] purchase_biz_no		-- 매입처 사업자등록번호

	//인자 확인
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//Bundle INVOKE
	var bundle = Bundle{Id_no: args[1], Barcode_id: args[2], Package_date: args[3], Part: args[4], Weight: args[5], Purchase_nm: args[6], Purchase_biz_no: args[7]}

	bundleAsBytes, _ := json.Marshal(bundle)
	jsonString := string(bundleAsBytes)
	log.Println("Logging: " + jsonString)
	//Bundle Asset 등록
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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[1], cowAsBytes)

	return shim.Success(nil)
}

//묶음번호 등록 - 자산(판매장)
func (s *SmartContract) registerInSalesBundleNum(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//판매장묶음번호등록

	log.Println("--==registerInSalesBundleNum==--")

	//'{'"Args":["registerInSalesBundleNum","BUNDLE0", "COW10", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- 개체식별번호
	//args[2] barcode_id			-- 바코드 ID
	//args[3] package_date			-- 포장처리일
	//args[4] part					-- 부위
	//args[5] weight				-- 중량
	//args[6] purchase_nm			-- 매입/의뢰처 상호
	//args[7] purchase_biz_no		-- 매입처 사업자등록번호

	//인자 확인
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//Bundle INVOKE
	var bundle = Bundle{Id_no: args[1], Barcode_id: args[2], Package_date: args[3], Part: args[4], Weight: args[5], Purchase_nm: args[6], Purchase_biz_no: args[7]}

	bundleAsBytes, _ := json.Marshal(bundle)
	jsonString := string(bundleAsBytes)
	log.Println("Logging: " + jsonString)
	//Bundle Asset 등록
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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[1], cowAsBytes)

	return shim.Success(nil)
}

//가축 이동(양수/양도), 도축장 이동, 가공장 이동, 판매장 이동
func (s *SmartContract) changeCowOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{"Args":["changeCowOwner", "COW3", "OWNER1", "OWNER2"]}'
	//args[0]				-- 소아이디
	//args[1]				-- 양수할소유자아이디
	//args[2]				-- 양도될소유자아이디

	log.Println("--==changeCowOwner==--")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}
	//OWNER 1,2,3에 적용
	json.Unmarshal(cowAsBytes, &cow)

	// if cow.Owner.Owner_id != args[1] {
	// 	return shim.Error("양도할 소의 소유자 아이디가 올바르지 않습니다.")
	// }

	//양도받을 소유자 정보 가져오기
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

//친환경 식품인증등록
func (s *SmartContract) addAut(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addAut==--")

	//'{'"Args":["addAut","OWNER90", "AutCheck", "201905251610", "축림", "1985.10.27", "전라북도 익산시", "Jeon Buk", "Cow", "30", "축평원", "10001", "20180525"]}'
	//aut_falg 			-- 인증구분
	//validity_date		-- 유효기간
	//farm_nm			-- 생상지명
	//farm_birth_date	-- 생산자생년월일
	//farm_addr			-- 생산자 주소
	//biz_addr			-- 사업장 소재지
	//aut_item			-- 인증품목
	//breed_head		-- 사육두수
	//aut_com			-- 인증기관
	//aut_id			-- 인증번호
	//aut_date			-- 인증일자

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

	//Json을 다시 Byte 형태로 변경
	ownerAsBytes, _ = json.Marshal(owner)
	//PutState해줌
	APIstub.PutState(args[0], ownerAsBytes)

	return shim.Success(nil)
}

//결핵/브루셀라(Tuberculousis/Brucella) 검사 - 거래
func (s *SmartContract) addBTVaccine(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addBTVaccine==--")

	//'{'"Args":["addBTVaccine","COW90", "farm_id", "farm_nm", "farm_addr", "farm_user_nm", "farm_user_birth", "farm_user_addr", "inspection_date", "inspection_head", "inspection_method", "livestock", "kind", "sex", "age", "id_no", "inspection_result", "inspection_part", "inspection_user_nm"]}'
	//args[0]						-- 소아이디
	//args[1] 	farm_id				-- 농장식별번호
	//args[2] 	farm_nm				-- 농장명
	//args[3] 	farm_addr			-- 소재지
	//args[4] 	farm_user_nm		-- 농장관리자 성명
	//args[5] 	farm_user_birth		-- 농장관리자 생년월일
	//args[6] 	farm_user_addr		-- 농장관리자 주소
	//args[7] 	inspection_date		-- 검사연월일
	//args[8] 	inspection_head		-- 검사두수
	//args[9] 	inspection_method	-- 검사방법
	//args[10]	livestock			-- 축종
	//args[11] 	kind				-- 품종
	//args[12] 	sex					-- 성별
	//args[13] 	age					-- 연령
	//args[14] 	id_no				-- 개체식별번호
	//args[15] 	inspection_result	-- 검사결과
	//args[16] 	inspection_part		-- 검사자소속
	//args[17] 	inspection_user_nm	-- 검사자성명

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//구제역 예방접종(Foot And Mouse Disease)
func (s *SmartContract) addFAMDVaccine(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addFAMDVaccine==--")

	//'{'"Args":["addFAMDVaccine","COW90", "farm_id", "farm_addr", "farm_tel", "breed_head", "item", "sex", "age", "id_no", "vaccination_date"]}'
	//args[0]					-- 소아이디
	//args[1] farm_id			-- 농장식별번호
	//args[2] farm_addr			-- 농장주소
	//args[3] farm_tel			-- 농장전화번호
	//args[4] breed_head		-- 사육두수
	//args[5] item				-- 품종
	//args[6] sex				-- 성별
	//args[7] age				-- 연령
	//args[8] id_no				-- 개체식별번호
	//args[9] vaccination_date	-- 예방접종일

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//가축 폐사 - 거래
func (s *SmartContract) addInfoDead(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoDead==--")

	//'{'"Args":["addInfoDead","COW90", "01", "180501-1", "180528", "Cancer", "burning"]}'
	//args[0]				-- 소아이디
	//args[1] farm_id		-- 농장아이디
	//args[2] id_no			-- 개체식별번호
	//args[3] det_date		-- 폐사일
	//args[4] det_reason	-- 폐사원안
	//args[5] det_method	-- 폐사처리방법

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//도축출하 - 거래
func (s *SmartContract) addInfoDeliver(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoDeliver==--")

	//'{'"Args":["addInfoDeliver","COW10", "180501-1", "RFID0"]}'
	//args[0]				-- Cow Key
	//args[1] id_no			-- 개체식별번호
	//args[2] rfid_no		-- RFID식별번호

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//도축검사등록 - 거래
func (s *SmartContract) addInfoInspect(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoInspect==--")

	//'{'"Args":["addInfoInspect","COW10", "COW", "180501-1", "300kg", "DoChuk1", "seal_10", "20180529", "FARM0", "Jeonju", "HACCP10", "Discard", "20180529", "Korea Inspect Center", "vetrinarian_100"]}'
	//args[0]						-- Cow Key
	//args[1] livestock				-- 가축의 종류
	//args[2] id_no					-- 개체식별번호
	//args[3] weight				-- 중량(지육 및 내장)
	//args[4] slaughter_nm			-- 도축장명
	//args[5] seal_no				-- 검인번호
	//args[6] slaughter_date		-- 도축연월일
	//args[7] farm_id				-- 농장식별번호
	//args[8] farm_addr				-- 도축의뢰인 주소
	//args[9] haccp_yn				-- HACCP 인증여부
	//args[10] fale_method			-- 불합격 식육 처분 방법
	//args[11] inspection_date		-- 검사연월일
	//args[12] inspection_part		-- 검사관 소속
	//args[13] inspection_user_nm	-- 검사관 성명
	//args[14] veterinarian_no		-- 수의사 면허번호

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//등급판정결과등록 - 거래
func (s *SmartContract) addInfoGradeResult(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoGradeResult==--")

	//'{'"Args":["addInfoGradeResult","COW10", "180529", "quality_part", "quality_nm", "subscriber_nm", "subscriber_birth", "subscriber_company", "subscriber_addr", "slaughter_nm", "slaughter_addr", "id_no", "weight", "meat_quality_grade", "meat_weight_grade", "grade_head"]}'
	//args[0]						-- Cow Key
	//args[1] grade_date			-- 등급판정연월일
	//args[2] quality_part			-- 품질평가사 소속
	//args[3] quality_nm			-- 품질평가사 성명
	//args[4] subscriber_nm			-- 신청인 성명
	//args[5] subscriber_birth		-- 신청인 생년월일
	//args[6] subscriber_company	-- 신청일 업소명
	//args[7] subscriber_addr		-- 신청인 주소
	//args[8] slaughter_nm			-- 도축장 명
	//args[9] slaughter_addr		-- 도축장 주소
	//args[10] id_no				-- 개체식별번호
	//args[11] weight				-- 중량(지육 및 내장)
	//args[12] meat_quality_grade	-- 육질등급
	//args[13] meat_weight_grade	-- 육량등급
	//args[14] grade_head			-- 판정두수

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//매입신고 - 거래(가공장)
func (s *SmartContract) addInfoInProcessesReportPurchase(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//가공장매입신고

	log.Println("--==addInfoInProcessesReportPurchase==--")

	//'{'"Args":["addInfoInProcessesReportPurchase","COW10", "barcode_id", "deal_date", "origin", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] barcode_id			-- 바코드ID
	//args[2] deal_date				-- 거래연월일
	//args[3] origin				-- 원산지
	//args[4] part					-- 부위
	//args[5] weight				-- 중량
	//args[6] purchase_nm			-- 매입처 상호
	//args[7] purchase_biz_no		-- 매입처 사업자등록번호

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)
	return shim.Success(nil)
}

//포장처리실적신고 - 거래
func (s *SmartContract) addInfoReportPacking(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoReportPacking==--")

	//'{'"Args":["addInfoReportPacking","COW10", "id_no", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- 개체식별번호
	//args[2] barcode_id			-- 바코드 ID
	//args[3] package_date			-- 포장처리일
	//args[4] part					-- 부위
	//args[5] weight				-- 중량
	//args[6] purchase_nm			-- 매입/의뢰처 상호
	//args[7] purchase_biz_no		-- 매입처 사업자등록번호

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)
	return shim.Success(nil)
}

//판매신고 - 거래
func (s *SmartContract) addInfoReportSale(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoReportSale==--")

	//'{'"Args":["addInfoReportSale","COW10", "id_no", "barcode_id", "package_date", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] id_no					-- 개체식별번호
	//args[2] barcode_id			-- 바코드 ID
	//args[3] sale_date				-- 판매연월일
	//args[4] part					-- 부위
	//args[5] weight				-- 판매중량
	//args[6] sale_nm				-- 판매처 상호
	//args[7] sale_biz_no			-- 판매처 사업자등록번호

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//매입신고 - 거래(판매장)
func (s *SmartContract) addInfoInSalesReportPurchase(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	log.Println("--==addInfoInSalesReportPurchase==--")

	//'{'"Args":["addInfoInSalesReportPurchase","COW10", "barcode_id", "deal_date", "origin", "part", "weight", "purchase_nm", "purchase_biz_no"]}'
	//args[0]						-- Cow Key
	//args[1] barcode_id			-- 바코드ID
	//args[2] deal_date				-- 거래연월일
	//args[3] origin				-- 원산지
	//args[4] part					-- 부위
	//args[5] weight				-- 중량
	//args[6] purchase_nm			-- 매입처 상호
	//args[7] purchase_biz_no		-- 매입처 사업자등록번호

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

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)
	return shim.Success(nil)
}

//소 정보 삭제
func (s *SmartContract) deleteCow(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var jsonResp string
	var CowJSON Cow
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	cowId := args[0]

	// to maintain the make~grade index, we need to read the cow first and get 'Make' values.
	// make ~ grade 인덱스를 유지하려면 먼저 소의 정보를 읽고 '가공업체' 정보를 가져와야합니다.
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

//test(Remark 적용 Test)
func (s *SmartContract) addRemark(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	//'{'"Args":["addRemark","COW4", "addVaccine", "True"]}

	log.Println("--==addRemark==--")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	log.Println("args[0]: " + args[0] + "  args[1]: " + args[1] + "  args[2]: " + args[2])
	cowAsBytes, _ := APIstub.GetState(args[0])
	cow := Cow{}

	//Byte형태의 json을 기존 json으로 변경
	json.Unmarshal(cowAsBytes, &cow)

	//새 비고값으로 넣을 Key Value를 넣음
	remark1 := Remark{Key: args[1], Value: args[2]}

	//비고값을 넣을 데이터 선언
	//Remarks := []Remark{}

	cow.setCowRemark(remark1)

	// if cow.Remarks == nil {
	// 	cow.Remarks[0] = remark1
	// } else {
	// 	//기존의 Remarks 배열에 새 비고값을 추가
	// 	Remarks = append(cow.Remarks, remark1)
	// }
	log.Println("change [cow] cow.Remarks Key: " + cow.Remarks[0].Key + " Value: " + cow.Remarks[0].Value)

	//이전의 Remark 대신 추가된 Remarks로 다시 넣어줌
	//cow.Remarks = Remarks

	//Json을 다시 Byte 형태로 변경
	cowAsBytes, _ = json.Marshal(cow)
	//PutState해줌
	APIstub.PutState(args[0], cowAsBytes)

	return shim.Success(nil)
}

//test(Cow Remark 적용 Test)
func (cow *Cow) setCowRemark(remark1 Remark) []Remark {
	cow.Remarks = append(cow.Remarks, remark1)
	return cow.Remarks
}

//test(Owner Remark 적용 Test)
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
