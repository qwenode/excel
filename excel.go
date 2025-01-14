package excel

import (
	"log"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

//Excel ..
type Excel struct {
	obj *ole.IDispatch
}

func CoInit() {
	err := ole.CoInitializeEx(0,0)
	if err!= nil {
		log.Println("ole初始化失败:",err)
	}
}
func CoUnInit() {
	ole.CoUninitialize()
}
//NewExcel ..
func NewExcel() (*Excel,error) {

	unknown, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		log.Fatalln("in func NewExcel1:", err)
		return nil,err
	}
	obj, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Fatalln("in func NewExcel2:", err)
		return nil,err
	}
	e := new(Excel)
	e.obj = obj
	if _, err := oleutil.PutProperty(obj, "Visible", false); err != nil {
		log.Println("in func NewExcel3:", err)
		e.Close()
		return nil,err
	}
	if _, err := oleutil.PutProperty(obj, "DisplayAlerts", false); err != nil {
		log.Println("in func NewExcel4:", err)
		e.Close()
		return nil,err
	}

	return e,nil
}

func (e *Excel) workbooks() *ole.IDispatch {
	obj, err := e.obj.GetProperty("Workbooks")
	if err != nil {
		log.Println(err)
		e.Close()
	}
	return obj.ToIDispatch()
}

//NewWorkBook ..
func (e *Excel) NewWorkBook() *Workbook {
	obj, err := e.workbooks().CallMethod("Add")
	if err != nil {
		log.Println("in func NewWorkBook:", err)
		return nil
	}
	workbook := new(Workbook)
	workbook.obj = obj.ToIDispatch()

	return workbook
}

//OpenWorkBook ..
func (e *Excel) OpenWorkBook(file string) *Workbook {
	path, err := filepath.Abs(file)
	if err != nil {
		log.Println("in func OpenWorkBook:", err)
		return nil
	}

	obj, err := e.workbooks().CallMethod("Open", path)
	if err != nil {
		log.Println("in func OpenWorkBook:", err)
		return nil
	}
	workbook := new(Workbook)
	workbook.obj = obj.ToIDispatch()

	return workbook
}

//Close ..
func (e *Excel) Close() {
	wbs := e.WorkBooks()
	for i := range wbs {
		wbs[i].Close()
	}

	e.obj.CallMethod("Quit")
	e.obj.Release()
}

//Visible ..
func (e *Excel) Visible(v bool) {
	if _, err := oleutil.PutProperty(e.obj, "Visible", v); err != nil {
		log.Println("in func Visible:", err)
	}
}

//Alert ..
func (e *Excel) Alert(v bool) {
	if _, err := oleutil.PutProperty(e.obj, "DisplayAlerts", v); err != nil {
		log.Println("in func Alert:", err)
	}
}

//WorkBooks ..
func (e *Excel) WorkBooks() []*Workbook {
	workbooksObj := e.workbooks()
	obj, err := workbooksObj.GetProperty("Count")
	if err != nil {
		log.Println("in func WorkBooks:", err)
	}
	len := int(obj.Val)

	workbooks := make([]*Workbook, 0)
	for i := 0; i < len; i++ {
		if obj, err := workbooksObj.GetProperty("Item", i+1); err != nil {
			log.Println("in func WorkBooks:", err)
		} else {
			workbook := new(Workbook)
			workbook.obj = obj.ToIDispatch()
			workbooks = append(workbooks, workbook)
		}

	}

	return workbooks
}
