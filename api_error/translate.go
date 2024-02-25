package api_error

type I18n struct {
	Ar string
	En string
}

var translated_errors = map[string]I18n{
	"something_went_wrong":                {Ar: "حدث خطأ ما", En: "Something went wrong"},
	"not_allowed_to_perform_this_request": {Ar: "غير مسموح لك بتنفيذ هذا الطلب", En: "You are not allowed to perform this request"},
	"resource_not_found":                  {Ar: "لم يتم العثور على المورد", En: "Resource not found"},
	"missing_authentication_data":         {Ar: "بيانات المصادقة مفقودة", En: "Missing authentication data"},
	"invalid_registration_data":           {Ar: "بيانات التسجيل غير صالحة", En: "Invalid registration data"},
	"invalid_login_data":                  {Ar: "بيانات تسجيل الدخول غير صالحة", En: "Invalid login data"},
	"missing_both_email_and_username":     {Ar: "يفتقد كل من البريد الإلكتروني واسم المستخدم", En: "Missing both email and username"},
	"invalid_email":                       {Ar: "بريد إلكتروني غير صالح", En: "Invalid email"},
	"user_not_found":                      {Ar: "المستخدم غير موجود", En: "User not found"},
	"old_password_do_not_match":           {Ar: "", En: ""},
	"image_type_png_only":                 {Ar: "", En: ""},
}

func (i I18n) Translate(lang string) string {
	if lang == "ar" {
		return i.Ar
	}
	return i.En
}
