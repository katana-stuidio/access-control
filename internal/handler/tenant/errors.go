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

// Erros Message Here

var ErroHttpMsgTenantNameIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Tenant Name is required",
	Code: http.StatusBadRequest,
}
var ErroHttpMsgSchemaNameIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Name is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantAlreadyExist handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Tenant Already Exist",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenanIsActiveIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Tenant IsActive is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantNotFound handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Tenant Not Found",
	Code: http.StatusNotFound,
}

var ErroHttpMsgToParseRequestTenantToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to parse Request Tenant to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToParseResponseTenantToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to parse Response Tenant to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToConvertingResponseTenantListToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to converting Response Tenant List to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToInsertTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to Insert the Tenant",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToUpdateTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to Update the Tenant",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToDeleteTenant handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to Delete the Tenant",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgTenantIdIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Tenant ID is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgTenantSchemaNameIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Tenant Schema Name is required",
	Code: http.StatusBadRequest,
}
