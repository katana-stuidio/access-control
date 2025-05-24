package tenant

import (
	"net/http"

	"github.com/katana-stuidio/access-control/internal/handler"
)

// Success Message Here
var SuccessHttpMsgToDeleteTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Ok Tenant Deleted",
	Code: http.StatusOK,
}

// Errors Message Here
var ErroHttpMsgTenantIdIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error Tenant ID is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantNameIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error Tenant Name is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantCNPJIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error Tenant CNPJ is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantCNPJIsInvalid handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error Tenant CNPJ is invalid",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantAlreadyExist handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error Tenant Already Exists",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantNotFound handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error Tenant Not Found",
	Code: http.StatusNotFound,
}

var ErroHttpMsgToParseRequestTenantToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error to parse Request Tenant to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToParseResponseTenantToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error to parse Response Tenant to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToConvertingResponseTenantListToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error to converting Response Tenant List to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToInsertTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error to Insert the Tenant",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToUpdateTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error to Update the Tenant",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToDeleteTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Error to Delete the Tenant",
	Code: http.StatusInternalServerError,
}
