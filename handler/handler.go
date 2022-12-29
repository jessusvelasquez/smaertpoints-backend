package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smartpoints/models"
	"smartpoints/services"

	"github.com/aws/aws-lambda-go/events"
)

func StartHandler(event events.APIGatewayProxyRequest) (interface{}, error) {
	fmt.Println("Start process with event: ", event.Path)
	fmt.Println("Start process with body: ", event.QueryStringParameters)

	response := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,HEAD,POST,OPTIONS",
		},
	}

	var inputData models.InputPostData
	switch action := event.HTTPMethod; action {
	case "POST":
		fmt.Println("Save data")
		err := json.Unmarshal([]byte(event.Body), &inputData)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = fmt.Sprintf(`{"message:":"%s"}`, "Error al procesar la información")
			break
		}
		err = services.SavePoints(inputData.ClientID, inputData.Puntos)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = fmt.Sprintf(`{"message":"%s"}`, "Error al guardar la información")
			break
		}
		response.StatusCode = http.StatusOK
		response.Body = fmt.Sprintf(`{"message":"%s"}`, "Información almacenada con éxito")
	case "GET":
		clientId := event.QueryStringParameters

		id := clientId["id"]

		item, err := services.GetPoints(id)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = fmt.Sprintf(`{"message":"%s"}`, "Cliente no encontrado")
			break
		}
		if item.ClientID == 0 {
			response.StatusCode = http.StatusNotFound
			response.Body = fmt.Sprintf(`{"message":"%s"}`, "Cliente no encontrado")
			break
		}

		response.StatusCode = http.StatusOK
		response.Body = fmt.Sprintf(`{"message":"%d"}`, item.Points)
	case "PUT":
		fmt.Println("Update data")
		err := json.Unmarshal([]byte(event.Body), &inputData)
		if err != nil {
			response.StatusCode = http.StatusInternalServerError
			response.Body = fmt.Sprintf(`{"message:":"%s"}`, "Error al procesar la información")
			break
		}
		err = services.UpdatePoints(inputData.ClientID, inputData.Puntos, event.Path)
		if err != nil {
			response.StatusCode = http.StatusNotFound
			response.Body = fmt.Sprintf(`{"message:":"%s"}`, "Error al guardar la información")
			break
		}
		response.StatusCode = http.StatusOK
		response.Body = fmt.Sprintf(`{"message:":"%s"}`, "Puntos actualizados")
	default:
		fmt.Println("No action valid")
		response.StatusCode = http.StatusBadRequest
		response.Body = fmt.Sprintf(`{"message":"%s"}`, "Solicitud no realizada")
	}

	return response, nil
}
