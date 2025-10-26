package handlers_v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/go-playground/validator/v10"

	"api_fiber/src/database"
	validators_common "api_fiber/src/validators"
	validators_request "api_fiber/src/validators/request"
	validators_response "api_fiber/src/validators/response"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// checkExistsUserWithPassword проверяет существование пользователя и валидность пароля.
//
// Параметры:
//   - db: указатель на подключение к БД
//   - checkUsername: имя пользователя для проверки
//   - checkPassword: пароль для проверки
//
// Возвращаемые коды статуса:
//   - 200: пользователь существует и пароль верный
//   - 401: неверный пароль
//   - 404: пользователь не найден
//   - 500: внутренняя ошибка сервера
func checkExistsUserWithPassword(db *database.DB, checkUsername string, checkPassword string) (int, *validators_response.ErrorResponse) {
	sql, args, err := database.Psql.
		Select("password").From("users").
		Where("username = ?", checkUsername).
		Where("disabled = false").
		ToSql()

	if err != nil {
		return 500, &validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query for check user",
			Details: err.Error(),
		}
	}

	var storedPassword string
	err = db.Connection.QueryRow(sql, args...).Scan(&storedPassword)

	if err != nil {
		return 500, &validators_response.ErrorResponse{
			Error:   "Database error",
			Details: err.Error(),
		}
	}

	if storedPassword == "" {
		return 404, &validators_response.ErrorResponse{
			Error:   "User not found",
			Details: "",
		}
	}

	// INFO! представим, что тут  как-будто сделал преобразование и проверил хэши
	if checkPassword != storedPassword {
		return 401, &validators_response.ErrorResponse{
			Error:   "Permission denied",
			Details: "Invalid password",
		}
	}

	return 200, nil
}

// CreateUser создает нового пользователя в системе.
//
// Тело запроса: RegisterRequest.
//
// Коды статуса ответа:
//   - 201: пользователь успешно создан
//   - 400: невалидные данные запроса
//   - 500: внутренняя ошибка сервера
func CreateUser(c *fiber.Ctx) error {
	var registerRequest validators_request.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := validate.Struct(registerRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
	}

	db := database.GetDB()

	sql, args, err := database.Psql.
		Insert("users").Columns("username", "password").
		Values(registerRequest.Username, registerRequest.Password).
		ToSql()

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query",
			Details: err.Error(),
		})
	}

	// INFO! допустипум пароль хешируется перед сохранением в бд
	_, err = db.Connection.Exec(sql, args...)

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Database error",
			Details: err.Error(),
		})
	}

	return c.Status(201).JSON(validators_response.EmptySuccessfulResponse{})
}

// UpdateCommonDataUser обновляет общие данные пользователя.
//
// Тело запроса: CommonDataUpdateRequest.
//
// Коды статуса ответа:
//   - 200: данные успешно обновлены
//   - 400: невалидные данные запроса
//   - 401: неверный пароль
//   - 404: пользователь не найден
//   - 500: внутренняя ошибка сервера
func UpdateCommonDataUser(c *fiber.Ctx) error {
	var commonDataUpdateRequest validators_request.CommonDataUpdateRequest
	if err := c.BodyParser(&commonDataUpdateRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := validate.Struct(commonDataUpdateRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
	}

	db := database.GetDB()

	statusCode, errorResponce := checkExistsUserWithPassword(db, commonDataUpdateRequest.Username, commonDataUpdateRequest.CurrentPassword)
	if errorResponce != nil {
		return c.Status(statusCode).JSON(errorResponce)
	}

	sql, args, err := database.Psql.
		Update("users").
		Set("about", commonDataUpdateRequest.About).
		Set("age", commonDataUpdateRequest.Age).
		Where("username = ?", commonDataUpdateRequest.Username).
		Where("disabled = false").
		ToSql()

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query for update user",
			Details: err.Error(),
		})
	}

	_, err = db.Connection.Exec(sql, args...)

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Database error",
			Details: err.Error(),
		})
	}

	return c.Status(200).JSON(validators_response.EmptySuccessfulResponse{})
}

// UpdatePasswordUser обновляет пароль пользователя.
//
// Тело запроса: PasswordUpdateRequest.
//
// Коды статуса ответа:
//   - 200: пароль успешно обновлен
//   - 400: невалидные данные запроса
//   - 401: неверный текущий пароль
//   - 404: пользователь не найден
//   - 500: внутренняя ошибка сервера
func UpdatePasswordUser(c *fiber.Ctx) error {
	var passwordUpdateRequest validators_request.PasswordUpdateRequest
	if err := c.BodyParser(&passwordUpdateRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := validate.Struct(passwordUpdateRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
	}

	db := database.GetDB()

	statusCode, errorResponce := checkExistsUserWithPassword(db, passwordUpdateRequest.Username, passwordUpdateRequest.CurrentPassword)
	if errorResponce != nil {
		return c.Status(statusCode).JSON(errorResponce)
	}

	sql, args, err := database.Psql.
		Update("users").
		Set("password", passwordUpdateRequest.NewPassword).
		Where("username = ?", passwordUpdateRequest.Username).
		Where("disabled = false").
		ToSql()

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query for update password user",
			Details: err.Error(),
		})
	}

	_, err = db.Connection.Exec(sql, args...)

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Database error",
			Details: err.Error(),
		})
	}

	return c.Status(200).JSON(validators_response.EmptySuccessfulResponse{})
}

// DisableUser отключает пользователя (помечает как disabled и далее он как-будто удален).
//
// Тело запроса: DisabledUserRequest.
//
// Коды статуса ответа:
//   - 204: пользователь успешно отключен
//   - 400: невалидные данные запроса
//   - 401: неверный пароль
//   - 404: пользователь не найден
//   - 500: внутренняя ошибка сервера
func DisableUser(c *fiber.Ctx) error {
	var disabledUserRequest validators_request.DisabledUserRequest
	if err := c.BodyParser(&disabledUserRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := validate.Struct(disabledUserRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
	}

	db := database.GetDB()

	statusCode, errorResponce := checkExistsUserWithPassword(db, disabledUserRequest.Username, disabledUserRequest.Password)
	if errorResponce != nil {
		return c.Status(statusCode).JSON(errorResponce)
	}

	sql, args, err := database.Psql.
		Update("users").
		Set("disabled", true).
		Where("username = ?", disabledUserRequest.Username).
		Where("disabled = false").
		ToSql()

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query for update password user",
			Details: err.Error(),
		})
	}

	_, err = db.Connection.Exec(sql, args...)

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Database error",
			Details: err.Error(),
		})
	}

	return c.Status(204).JSON(validators_response.EmptySuccessfulResponse{})
}

// ListUsers возвращает список пользователей с пагинацией.
//
// Тело запроса: UsersPageRequest.
//
// Коды статуса ответа:
//   - 200: список пользователей успешно получен
//   - 400: невалидные параметры пагинации
//   - 500: внутренняя ошибка сервера
func ListUsers(c *fiber.Ctx) error {
	var usersPageRequest validators_request.UsersPageRequest
	if err := c.BodyParser(&usersPageRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
	}

	if err := validate.Struct(usersPageRequest); err != nil {
		return c.Status(400).JSON(validators_response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
	}

	db := database.GetDB()

	sql, args, err := database.Psql.
		Select("username", "about", "age", "COUNT(1) OVER() as total_count").
		From("users").
		Where("disabled = false").
		Limit(usersPageRequest.SizePage).
		Offset((usersPageRequest.Page - 1) * usersPageRequest.SizePage).
		ToSql()

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query for get list users",
			Details: err.Error(),
		})
	}

	var listUsers validators_response.ListUsers = validators_response.ListUsers{
		Users:       make([]validators_common.User, 0),
		SizePage:    usersPageRequest.SizePage,
		CurrentPage: usersPageRequest.Page,
	}

	rows, err := db.Connection.Query(sql, args...)
	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Database error",
			Details: err.Error(),
		})
	}
	defer rows.Close()

	totalCount := uint64(0)

	for rows.Next() {
		var user validators_common.User
		err := rows.Scan(&user.Username, &user.About, &user.Age, &totalCount)

		if err != nil {
			return c.Status(500).JSON(validators_response.ErrorResponse{
				Error:   "Error scanning user data",
				Details: err.Error(),
			})
		}

		listUsers.Users = append(listUsers.Users, user)
	}

	listUsers.AmountPage = (totalCount + usersPageRequest.SizePage - 1) / usersPageRequest.SizePage

	return c.Status(200).JSON(listUsers)
}

// GetUser возвращает данные конкретного пользователя по имени.
//
// Параметры пути:
//   - username: имя пользователя
//
// Коды статуса ответа:
//   - 200: данные пользователя успешно получены
//   - 404: пользователь не найден
//   - 500: внутренняя ошибка сервера
func GetUser(c *fiber.Ctx) error {
	usernameParameter := c.Params("username")

	db := database.GetDB()

	sql, args, err := database.Psql.
		Select("username", "about", "age").
		From("users").
		Where("username = ?", usernameParameter).
		Where("disabled = false").
		ToSql()

	if err != nil {
		return c.Status(500).JSON(validators_response.ErrorResponse{
			Error:   "Unable to assemble SQL query for get list users",
			Details: err.Error(),
		})
	}

	var user validators_common.User
	err = db.Connection.QueryRow(sql, args...).Scan(&user.Username, &user.About, &user.Age)

	if err != nil {
		return c.Status(404).JSON(validators_response.ErrorResponse{
			Error:   "User not found",
			Details: "The query returned zero rows",
		})
	}

	return c.Status(200).JSON(user)
}
