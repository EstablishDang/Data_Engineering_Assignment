package service

import (
	proto_mssql_server "app_mssql_server/mssql_server/proto/v1"
	"app_mssql_server/utils"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	//log logger.Logger
	proto_mssql_server.AppMssqlMgmtServiceServer
}

// Replace with your own connection parameters
var server = "3.0.135.17"
var port = 1433
var user = "admin"
var password = "lhp21051999_phat"
var database = "AdventureWorks2008R2"

// var server = "localhost"
// var port = 1433
// var user = "sa"
// var password = "Lap20901359$"
// var database = "AdventureWorks2008R2"

var db *sql.DB

func InitDB() {
	//Test
	// Create connection string
	var err error
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: " + err.Error())
	}
	log.Printf("Connected!\n")

}

func (s *Service) GetVendor(ctx context.Context, req *proto_mssql_server.GetVendorReq) (*proto_mssql_server.ListVendorResp, error) {
	var Vendor_Info []*proto_mssql_server.VendorInfo

	if req.Type == "All" {
		cmd := "select Name, CreditRating, ActiveFlag, PreferredVendorStatus, AccountNumber from Purchasing.Vendor;"
		rows, err := db.QueryContext(ctx, cmd)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var info proto_mssql_server.VendorInfo
			var name string
			var CreditRating uint32
			var PreferredVendorStatus bool
			var ActiveFlag bool
			var AccountNumber string

			err := rows.Scan(&name, &CreditRating, &ActiveFlag, &PreferredVendorStatus, &AccountNumber)
			if err != nil {
				return nil, err
			}
			if CreditRating == 1 {
				info.CreditRating = "Superior"
			} else if CreditRating == 2 {
				info.CreditRating = "Excellent"
			} else if CreditRating == 3 {
				info.CreditRating = "Above average"
			} else if CreditRating == 4 {
				info.CreditRating = "Average"
			} else if CreditRating == 5 {
				info.CreditRating = "Below average"
			}

			info.Name = name
			info.Active = ActiveFlag
			info.AccNumber = AccountNumber
			info.PreferVendor = PreferredVendorStatus

			Vendor_Info = append(Vendor_Info, &info)
		}
		return &proto_mssql_server.ListVendorResp{VendorInfo: Vendor_Info}, nil
	} else if req.Type == "Name" {
		cmd := "select Name from Purchasing.Vendor;"
		rows, err := db.QueryContext(ctx, cmd)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var info proto_mssql_server.VendorInfo
			var name string

			err := rows.Scan(&name)
			if err != nil {
				return nil, err
			}
			info.Name = name

			Vendor_Info = append(Vendor_Info, &info)
		}
		return &proto_mssql_server.ListVendorResp{VendorInfo: Vendor_Info}, nil
	} else if req.Type != "" {
		cmd := fmt.Sprintf("select AccountNumber, CreditRating, PreferredVendorStatus, ActiveFlag from Purchasing.Vendor WHERE Name = '%s';", req.Type)
		rows, err := db.QueryContext(ctx, cmd)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var info proto_mssql_server.VendorInfo
			var CreditRating uint32
			var PreferredVendorStatus bool
			var ActiveFlag bool
			var AccountNumber string

			err := rows.Scan(&AccountNumber, &CreditRating, &ActiveFlag, &PreferredVendorStatus)
			if err != nil {
				return nil, err
			}
			if CreditRating == 1 {
				info.CreditRating = "Superior"
			} else if CreditRating == 2 {
				info.CreditRating = "Excellent"
			} else if CreditRating == 3 {
				info.CreditRating = "Above average"
			} else if CreditRating == 4 {
				info.CreditRating = "Average"
			} else if CreditRating == 5 {
				info.CreditRating = "Below average"
			}

			info.Name = req.Type
			info.Active = ActiveFlag
			info.AccNumber = AccountNumber
			info.PreferVendor = PreferredVendorStatus
			Vendor_Info = append(Vendor_Info, &info)
		}
		if Vendor_Info == nil {
			return nil, status.Error(codes.NotFound, "Not found info")
		}
		return &proto_mssql_server.ListVendorResp{VendorInfo: Vendor_Info}, nil
	}

	return nil, status.Error(codes.NotFound, "Not found info")
}

func (s *Service) ListShipMethod(ctx context.Context, req *empty.Empty) (*proto_mssql_server.ListShipMethodResp, error) {
	var shipmethod []*proto_mssql_server.ShipMethod

	cmd := "select Name, ShipBase, ShipRate from Purchasing.ShipMethod;"

	rows, err := db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var info proto_mssql_server.ShipMethod
		var name string
		var ship_base float32
		var ship_rate float32

		err := rows.Scan(&name, &ship_base, &ship_rate)
		if err != nil {
			return nil, err
		}
		info.Name = name
		info.ShipBase = ship_base
		info.ShipRate = ship_rate

		shipmethod = append(shipmethod, &info)
	}

	return &proto_mssql_server.ListShipMethodResp{ShipMethod: shipmethod}, nil
}

func (s *Service) GetPurchaseOrderInfo(ctx context.Context, req *proto_mssql_server.PurchaseOrderInfoReq) (*proto_mssql_server.PurchaseOrderInfo, error) {
	var PurchaseHeader proto_mssql_server.PurchaseOrderHeader
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid Input ID of PurchaseOder")
	}

	cmd := fmt.Sprintf("select * from Purchasing.PurchaseOrderHeader WHERE PurchaseOrderID='%d';", req.Id)
	rows, err := db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&PurchaseHeader.PurchaseOrderID,
			&PurchaseHeader.RevisionNumber,
			&PurchaseHeader.Status,
			&PurchaseHeader.EmployeeID,
			&PurchaseHeader.VendorID,
			&PurchaseHeader.ShipMethodID,
			&PurchaseHeader.OrderDate,
			&PurchaseHeader.ShipDate,
			&PurchaseHeader.SubTotal,
			&PurchaseHeader.TaxAmt,
			&PurchaseHeader.Freight,
			&PurchaseHeader.TotalDue,
			&PurchaseHeader.ModifiedDate)
		if err != nil {
			return nil, err
		}
	}
	var status_info string //Description	Order current status. 1 = Pending; 2 = Approved; 3 = Rejected; 4 = Complete
	if PurchaseHeader.Status == 1 {
		status_info = "Pending"
	} else if PurchaseHeader.Status == 2 {
		status_info = "Approved"
	} else if PurchaseHeader.Status == 3 {
		status_info = "Rejected"
	} else if PurchaseHeader.Status == 4 {
		status_info = "Completed"
	}

	// look up Employee 1
	cmd = fmt.Sprintf("select JobTitle, Gender from HumanResources.Employee WHERE BusinessEntityID='%d';", PurchaseHeader.EmployeeID)
	rows, err = db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Employee_info proto_mssql_server.Employee
	for rows.Next() {
		err := rows.Scan(&Employee_info.JobTitle,
			&Employee_info.Gender)
		if err != nil {
			return nil, err
		}
	}

	// look up Employee 2
	var d_id, s_id uint32
	cmd = fmt.Sprintf("select DepartmentID, ShiftID from HumanResources.EmployeeDepartmentHistory WHERE BusinessEntityID='%d';", PurchaseHeader.EmployeeID)
	rows, err = db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&d_id, &s_id)
		if err != nil {
			return nil, err
		}
	}

	// look up Employee 3
	cmd = fmt.Sprintf("select Name, GroupName from HumanResources.Department WHERE DepartmentID='%d';", d_id)
	rows, err = db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&Employee_info.DepartmentName, &Employee_info.GroupName)
		if err != nil {
			return nil, err
		}
	}

	// look up Employee 4
	cmd = fmt.Sprintf("select Name from HumanResources.Shift WHERE ShiftID='%d';", s_id)
	rows, err = db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&Employee_info.Shift)
		if err != nil {
			return nil, err
		}
	}

	// look up verdor
	var credit_rating int16
	var vendor_info proto_mssql_server.VendorInfo
	cmd = fmt.Sprintf("select AccountNumber, Name, CreditRating, PreferredVendorStatus, ActiveFlag from Purchasing.Vendor WHERE BusinessEntityID = '%d';", PurchaseHeader.VendorID)
	rows, err = db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&vendor_info.AccNumber,
			&vendor_info.Name,
			&credit_rating,
			&vendor_info.PreferVendor,
			&vendor_info.Active)
		if err != nil {
			return nil, err
		}
	}
	var CreditRating string
	//Description	1 = Superior, 2 = Excellent, 3 = Above average, 4 = Average, 5 = Below average
	if credit_rating == 1 {
		CreditRating = "Superior"
	} else if credit_rating == 2 {
		CreditRating = "Excellent"
	} else if credit_rating == 3 {
		CreditRating = "Above average"
	} else if credit_rating == 4 {
		CreditRating = "Average"
	} else if credit_rating == 5 {
		CreditRating = "Below average"
	}

	// look up ship method
	var ship_method proto_mssql_server.ShipMethod
	cmd = fmt.Sprintf("select Name, ShipBase, ShipRate from Purchasing.ShipMethod WHERE ShipMethodID = '%d';", PurchaseHeader.ShipMethodID)
	rows, err = db.QueryContext(ctx, cmd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&ship_method.Name,
			&ship_method.ShipBase,
			&ship_method.ShipRate)
		if err != nil {
			return nil, err
		}
	}

	return &proto_mssql_server.PurchaseOrderInfo{
		EmployeeInfo: &proto_mssql_server.Employee{
			JobTitle:       Employee_info.JobTitle,
			Gender:         Employee_info.Gender,
			DepartmentName: Employee_info.DepartmentName,
			GroupName:      Employee_info.GroupName,
			Shift:          Employee_info.Shift,
		},
		VendorInfo: &proto_mssql_server.VendorInfo{
			AccNumber:    vendor_info.AccNumber,
			Name:         vendor_info.Name,
			CreditRating: CreditRating,
			PreferVendor: vendor_info.PreferVendor,
			Active:       vendor_info.Active,
		},
		ShipMethod: &proto_mssql_server.ShipMethod{
			Name:     ship_method.Name,
			ShipBase: ship_method.ShipBase,
			ShipRate: ship_method.ShipRate,
		},
		Status:       status_info,
		OrderDate:    PurchaseHeader.OrderDate,
		ShipDate:     PurchaseHeader.ShipDate,
		SubTotal:     PurchaseHeader.SubTotal,
		TaxAmt:       PurchaseHeader.TaxAmt,
		Freight:      PurchaseHeader.Freight,
		TotalDue:     PurchaseHeader.TotalDue,
		ModifiedDate: PurchaseHeader.ModifiedDate,
	}, nil
}

func (s *Service) AddNewVendor(ctx context.Context, req *proto_mssql_server.AddNewVendorReq) (*empty.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid Input of Vendor")
	}
	now := time.Now()
	time_stamp := now.Unix()
	time_string := utils.UnixSecondsToTimestampRFC3339(time_stamp)

	uuid, _ := NewUUID()
	// Create id
	var new_id uint16
	tsql := fmt.Sprintf("INSERT INTO Person.BusinessEntity (rowguid, ModifiedDate) VALUES ('%s', '%s'); SELECT SCOPE_IDENTITY()", uuid, time_string)
	err := db.QueryRow(tsql).Scan(&new_id)
	if err != nil {
		fmt.Println("Error inserting new row vendor: " + err.Error())
		return nil, err
	}
	fmt.Printf("new_ID: %d\n", new_id)
	if err != nil {
		return nil, err
	}

	// add new
	acc_name := strings.ToUpper(req.Name) + "0001"
	tsql = fmt.Sprintf("INSERT INTO Purchasing.Vendor (BusinessEntityID, AccountNumber, Name, CreditRating, PreferredVendorStatus, ActiveFlag,PurchasingWebServiceURL,ModifiedDate) VALUES ('%d', '%s', '%s', '%s', '%t', '%t', '%s', '%s');", new_id, acc_name, req.Name, "1", true, true, req.UrlWeb, time_string)
	_, err = db.Exec(tsql)
	if err != nil {
		fmt.Println("Error inserting new row vendor: " + err.Error())
		return nil, err
	}

	return &empty.Empty{}, nil
}
