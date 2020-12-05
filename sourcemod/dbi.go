/**
 * sourcemod/dbi.go
 * 
 * Copyright 2020 Nirari Technologies, Alliedmodders LLC.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 * 
 */

package main


type DBResult int
const (
	DBVal_Error = DBResult(0)        /**< Column number/field is invalid. */
	DBVal_TypeMismatch = 1 /**< You cannot retrieve this data with this type. */
	DBVal_Null = 2         /**< Field has no data (NULL) */
	DBVal_Data = 3         /**< Field has data */
)

type DBBindType int
const (
	DBBind_Int = DBResult(0)         /**< Bind an integer. */
	DBBind_Float = DBResult(1)       /**< Bind a float. */
	DBBind_String = DBResult(2)      /**< Bind a string. */
)

type DBPriority int
const (
	DBPrio_High = DBPriority(0)        /**< High priority. */
	DBPrio_Normal = DBPriority(1)      /**< Normal priority. */
	DBPrio_Low = DBPriority(2)         /**< Low priority. */
)

type DBDriver Handle
/**
 * func static Find(name string) DBDriver
 * __sp__(`DBDriver db = DBDriver.Find(name);`)
 */
func (DBDriver) GetIdentifier(ident []char, maxlength int)
func (DBDriver) GetProduct(ident []char, maxlength int)


type DBResultSet struct {
	RowCount, FieldCount, AffectedRows, InsertId int
	HasResults, MoreRows bool
}

func (DBResultSet) FetchMoreResults() bool
func (DBResultSet) FieldNumToName(field int, name []char, maxlength int)
func (DBResultSet) FieldNameToNum(name string, field *int) bool
func (DBResultSet) FetchRow() bool
func (DBResultSet) Rewind() bool
func (DBResultSet) FetchString(field int, name []char, maxlength int, result *DBResult) int
func (DBResultSet) FetchFloat(field int, result *DBResult) float
func (DBResultSet) FetchInt(field int, result *DBResult) int
func (DBResultSet) IsFieldNull(field int) bool
func (DBResultSet) FetchSize(field int) int


type (
	SQLTxnSuccess func(db Database, data any, numQueries int, results []DBResultSet, queryData []any)
	SQLTxnFailure func(db Database, data any, numQueries int, err string, failIndex int, queryData []any)
)

type Transaction Handle
/// public native Transaction();
/// __sp__(`transact = new Transaction();`)

func (t Transaction) AddQuery(query string, data any) int


type DBStatement Handle
func (DBStatement) BindInt(param, number int, signed bool)
func (DBStatement) BindFloat(param int, value float)
func (DBStatement) BindString(param int, value string, copy bool)


type (
	SQLConnectCallback func(db Database, err string, data any)
	SQLQueryCallback func(db Database, results DBResultSet, err string, data any)
	
	/// static func Connect(callback SQLConnectCallback, name="default" string, data=0 any)
	Database struct {
		Driver DBDriver
	}
)

func (Database) SetCharset(charset string) bool
func (Database) Escape(str string, buffer []char, maxlength int, written *int) bool
func (Database) Format(buffer string, maxlength int, format string, fmt_args ...any) int
func (Database) IsSameConnection(other Database) bool
func (Database) Query(callback SQLQueryCallback, query string, data any, prio DBPriority)
func (Database) Execute(txn Transaction, onSuccess SQLTxnSuccess, onError SQLTxnFailure, data any, prio DBPriority)

func SQL_Connect(confname string, persistent bool, err []char, maxlength int) Database
func SQL_DefConnect(err []char, maxlength int, persistent bool) Database
func SQL_ConnectCustom(keyvalues KeyValues, err []char, maxlength int, persistent bool) Database
func SQLite_UseDatabase(database string, err []char, maxlength int) Database

func SQL_CheckConfig(name string) bool
func SQL_GetDriver(name string) DBDriver
func SQL_ReadDriver(database Database, ident []char, ident_length int) DBDriver
func SQL_SetCharset(database Database, charset string) bool

func SQL_GetDriverIdent(driver DBDriver, ident []char, maxlength int)
func SQL_GetDriverProduct(driver DBDriver, product []char, maxlength int)
func SQL_GetAffectedRows(hndl any) int
func SQL_GetInsertId(hndl any) int
func SQL_GetError(hndl any, err []char, maxlength int) bool

func SQL_EscapeString(database Database, str string, buffer []char, maxlength int, written *int) bool
func SQL_FormatQuery(database Database, buffer string, maxlength int, format string, fmt_args ...any) int
func SQL_QuoteString(database Database, str string, buffer []char, maxlength int, written *int) bool
func SQL_FastQuery(database Database, query string, len int) bool
func SQL_Query(database Database, query string, len int) DBResultSet
func SQL_PrepareQuery(database Database, query string, err []char, maxlength int) DBStatement

func SQL_FetchMoreResults(query DBResultSet) bool
func SQL_HasResultSet(query any) bool
func SQL_GetRowCount(query any) int
func SQL_GetFieldCount(query any) int
func SQL_FieldNumToName(query any, field int, name []char, maxlength int)
func SQL_FieldNameToNum(query any, name string, field *int) bool
func SQL_FetchRow(query any) bool
func SQL_MoreRows(query any) bool
func SQL_Rewind(query any) bool

func SQL_FetchString(query any, field int, name []char, maxlength int, result *DBResult) int
func SQL_FetchFloat(query any, field int, result *DBResult) float
func SQL_FetchInt(query any, field int, result *DBResult) int
func SQL_IsFieldNull(hndl any, field int) bool
func SQL_FetchSize(query any, field int) int

func SQL_BindParamInt(statement DBStatement, param, number int, signed bool)
func SQL_BindParamFloat(statement DBStatement, param int, value float)
func SQL_BindParamString(statement DBStatement, param int, value string, copy bool)

func SQL_Execute(statement DBStatement) bool
func SQL_LockDatabase(database Database)
func SQL_UnlockDatabase(database Database)

type SQLTCallback func(owner, hndl Handle, err string, data any)
func SQL_IsSameConnection(hndl1, hndl2 Database) bool
func SQL_TConnect(callback SQLTCallback, name string, data any);
func SQL_TQuery(database Database, callback SQLTCallback, query string, data any, prio DBPriority)

func SQL_CreateTransaction() Transaction
func SQL_AddQuery(txn Transaction, query string, data any) int
func SQL_ExecuteTransaction(db Database, txn Transaction, onSuccess SQLTxnSuccess, onError SQLTxnFailure, data any, priority DBPriority)