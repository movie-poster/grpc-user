package repository_user

import (
	"grpc-user/internal/constant"
	"grpc-user/internal/domain/entity"
	objectvalue "grpc-user/internal/domain/object-value"
	irepository "grpc-user/internal/domain/repository/interface"
	s "grpc-user/internal/domain/service"
	pb "grpc-user/internal/infra/proto/user"
	"grpc-user/internal/utils"
	"net/http"

	"gorm.io/gorm"
)

type crud struct {
	DB    *gorm.DB
	brevo s.IBrevoSender
}

func UserRepository(DB *gorm.DB, brevoService s.IBrevoSender) irepository.IUserCrud {
	return &crud{
		DB:    DB,
		brevo: brevoService,
	}
}

func (u *crud) Insert(user *entity.User, originalPassword string) *objectvalue.ResponseValue {
	var userQuery *entity.User
	tx := u.DB.Begin()

	err := u.DB.Where("nickname = ? AND state = ?", user.NickName, constant.ActiveState).
		Or("nick_name= ?", user.NickName).
		First(&userQuery).Error
	if err != nil {
		utils.LogWarning("Registro no encontrado", "No se ha encontrado el registro", "Insert", user)
	}

	if userQuery.ID == constant.IDDoesNotExists {
		err = tx.Create(&user).Error
		if err != nil {
			tx.Rollback()
			message := utils.CheckErrorFromDB(err)
			utils.LogError("Error al guardar el registro", message, "Insert", http.StatusBadRequest, user)
			return objectvalue.BadResponseSingle(message)
		}

		err = u.brevo.SendToNewUser(*user, originalPassword)
		if err != nil {
			tx.Rollback()
			utils.LogError("Error al guardar el registro", "No fue posible enviar la confirmación de email", "Insert", http.StatusBadRequest, user)
			return objectvalue.BadResponseSingle("No fue posible enviar la confirmación de email")
		}

		tx.Commit()
		utils.LogSuccess("Registro guardado", "Insert", user)
		return &objectvalue.ResponseValue{
			Title:   "¡Creado exitosamente!",
			IsOk:    true,
			Message: "El usuario se ha creado",
			Status:  http.StatusCreated,
			Value:   u.MarshalResponse(user, constant.HidePassword),
		}
	}

	utils.LogError("Error al guardar el registro", "El usuario ya está creado en la misma empresa", "Insert", http.StatusBadRequest, user)
	return objectvalue.BadResponseSingle("El usuario ya está creado en la misma empresa")
}

func (u *crud) Delete(ID uint64) *objectvalue.ResponseValue {
	err := u.DB.Model(&entity.User{}).Where("id", ID).Update("state", constant.InactiveState).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al eliminar el registro", message, "Delete", http.StatusBadRequest, ID)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no existoso",
			IsOk:    false,
			Message: message,
			Status:  http.StatusBadRequest,
		}
	}
	utils.LogSuccess("Registro eliminado", "Delete", ID)
	return &objectvalue.ResponseValue{
		Title:   "Eliminado exitosamente!",
		IsOk:    true,
		Message: "Se eliminó correctamente",
		Status:  http.StatusOK,
	}
}

func (u *crud) Update(user *entity.User) *objectvalue.ResponseValue {
	var userQuery *entity.User

	err := u.DB.Model(&entity.User{}).
		Where("nickname = ? and state = ?", user.NickName, constant.ActiveState).
		First(&userQuery).Error
	if err != nil {
		utils.LogWarning("Registro no encontrado", "No se ha encontrado el registro", "Update", user)
	}

	// Significa que el usuario ya existe
	if userQuery.ID != constant.IDDoesNotExists && user.ID != userQuery.ID {
		utils.LogWarning("Registro encontrado", "El usuario ya está creado en la misma empresa", "Update", user)
		return objectvalue.BadResponseSingle("El usuario ya está creado en la misma empresa")
	}

	// No existe o significa que es el mismo registro
	if userQuery.ID == constant.IDDoesNotExists || user.ID == userQuery.ID {

		err := u.DB.Where("id", user.ID).Updates(&user).Error
		if err != nil {
			message := utils.CheckErrorFromDB(err)
			utils.LogError("Error al actualizar el registro", message, "Update", http.StatusBadRequest, user)
			return objectvalue.BadResponseSingle(message)
		}

		utils.LogSuccess("Actualizado con éxito", "Update", user)
		return &objectvalue.ResponseValue{
			Title:   "Editado exitosamente!",
			IsOk:    true,
			Message: "El usuario fue editado correctamente",
			Status:  http.StatusOK,
			Value:   u.MarshalResponse(user, constant.HidePassword),
		}
	}

	utils.LogError("Error al actualizar el registro", "Errores al hacer la modificación", "Update", http.StatusBadRequest, user)
	return objectvalue.BadResponseSingle("Errores al hacer la modificación")
}

func (u *crud) VerifyUser(nickname string) *objectvalue.ResponseValue {
	var user *entity.User

	result := u.DB.Where("nick_name", nickname).First(&user)
	if err := result.Error; err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al buscar por ID", message, "VerifyUser", http.StatusBadRequest, nickname)
		return &objectvalue.ResponseValue{
			Title:   "Proceso no existoso",
			IsOk:    false,
			Message: message,
			Status:  http.StatusBadRequest,
			Value:   &pb.User{},
		}
	}
	utils.LogSuccess("Registro encontrado", "VerifyUser", nickname)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Se ha encontrado el usuario con el ID",
		Status:  http.StatusCreated,
		Value:   u.MarshalResponse(user, constant.ShowPassword),
	}
}

func (u *crud) ResetPassword(reset entity.RecoverPassword) *objectvalue.ResponseValue {
	tx := u.DB.Begin()

	var user entity.User
	err := u.DB.Where("nick_name", reset.Nickname).
		First(&user).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("No se ha podido encontrar el nombre de usuario", message, "ResetPassword", http.StatusBadRequest)
		return objectvalue.BadResponseSingle("No se ha podido encontrar el nombre de usuario")
	}

	err = tx.Create(&reset).Error
	if err != nil {
		tx.Rollback()
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Error al registrar nuevo token", message, "ResetPassword", http.StatusBadRequest)
		return objectvalue.BadResponseSingle("Error al registrar nuevo token")
	}

	err = u.brevo.ResetPasswordWithUsername(reset.Token, reset.Nickname, user.Email, user.Name)
	if err != nil {
		tx.Rollback()
		utils.LogError("No fue posible enviar el email", err.Error(), "ResetPassword", http.StatusBadRequest)
		return objectvalue.BadResponseSingle("Error al registrar nuevo token")
	}

	tx.Commit()
	utils.LogSuccess("Registro encontrado", "ResetPassword", reset)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Te hemos enviado un email para que puedas continuar con el proceso de recuperación",
		Status:  http.StatusCreated,
	}
}

func (u *crud) CheckToken(token string) *objectvalue.ResponseValue {
	var counter int64

	err := u.DB.Model(&entity.RecoverPassword{}).
		Where("token = ? AND state = ?", token, constant.ActiveState).
		Count(&counter).Error
	if err != nil || counter == 0 {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("No se ha podido encontrar el token, no existe", message, "CheckToken", http.StatusBadRequest)
		return objectvalue.BadResponseSingle("El token es inválido o ha expirado. Te pedimos que solicites nuevamente un cambio de contraseña. Recuerda que tienes dos horas para poder realizar el proceso")
	}

	utils.LogSuccess("Registro encontrado", "CheckToken", token)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "El token si existe",
		Status:  http.StatusOK,
	}
}

func (u *crud) ChangePassword(token, password, nickName string) *objectvalue.ResponseValue {
	var recover entity.RecoverPassword
	err := u.DB.Where("token = ? AND state = ?", token, constant.ActiveState).
		First(&recover).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("No se ha podido encontrar el token, no existe", message, "ChangePassword", http.StatusBadRequest)
		return objectvalue.BadResponseSingle("El token es inválido o ha expirado. Te pedimos que solicites nuevamente un cambio de contraseña. Recuerda que tienes dos horas para poder realizar el proceso")
	}
	err = u.DB.Model(&entity.User{}).
		Where("nick_name", nickName).
		Update("password", password).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("No fue posible actualizar el usuario", message, "ChangePassword", http.StatusBadRequest)
		return objectvalue.BadResponseSingle("No ha sido posible actualizar sus datos. Contacte a soporte por favor.")
	}

	err = u.DB.Model(&entity.RecoverPassword{}).
		Where("token", token).
		Update("state", constant.InactiveState).Error
	if err != nil {
		message := utils.CheckErrorFromDB(err)
		utils.LogError("Errores al actualizar RecoverPassword", message, "ChangePassword", http.StatusBadRequest)
	}

	utils.LogSuccess("Contraseña actualizada exitosamente", "ChangePassword", token)
	return &objectvalue.ResponseValue{
		Title:   "¡Proceso exitoso!",
		IsOk:    true,
		Message: "Contraseña actualizada exitosamente",
		Status:  http.StatusOK,
	}
}

func (u *crud) MarshalResponse(user *entity.User, addPasswordField bool) *pb.User {
	userPB := &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		NickName:  user.NickName,
		IsAdmin:   user.IsAdmin,
		State:     user.State,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}

	if addPasswordField {
		userPB.Password = user.Password
	}

	return userPB
}
