package user

import (
	"net/http"

	"github.com/katana-stuidio/access-control/internal/handler"
)

// Success Message Here

var SuccessHttpMsgToDeleteUser handler.HttpMsg = handler.HttpMsg{
	Msg:  "Ok User Deleted",
	Code: http.StatusOK,
}

// Erros Message Here

var ErroHttpMsgUserIdIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User ID is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserNameIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User Name(CPF) is required",
	Code: http.StatusBadRequest,
}
var ErroHttpMsgNameIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Name is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserPasswordIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro Password is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserAlreadyExist handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User Already Exist",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserPriceIsRequired handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User Name is required",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserCpfIsInvalid handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User Cpf is invalid",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserCnpjIsInvalid handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User CNPJ is invalid",
	Code: http.StatusBadRequest,
}

var ErroHttpMsgUserNotFound handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro User Not Found",
	Code: http.StatusNotFound,
}

var ErroHttpMsgToParseRequestUserToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to parse Request User to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToParseResponseUserToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to parse Response User to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToConvertingResponseUserListToJson handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to converting Response User List to JSON",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToInsertUser handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to Insert the User",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToUpdateUser handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to Update the User",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgToDeleteUser handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro to Delete the User",
	Code: http.StatusInternalServerError,
}

var ErroHttpMsgCNPJNotFound handler.HttpMsg = handler.HttpMsg{
	Msg:  "Erro CNPJ Not Found",
	Code: http.StatusNotFound,
}

var ErroHttpMsgInvalidRole handler.HttpMsg = handler.HttpMsg{
	Msg:  "Invalid role. Valid roles are: Professor, Estudante, Instituicao, Admin",
	Code: http.StatusBadRequest,
}
