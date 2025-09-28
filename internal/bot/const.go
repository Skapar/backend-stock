package bot

// Button Key
const (
	STARTBUTTON             = "start"
	GETPAYMENTDETAILSBUTTON = "get_payment_details"
	ASKRECEIPTBUTTON        = "send_receipt"
)

// Button Text
const (
	GETPAYMENTDETAILSBUTTONTEXT = "Получение реквизитов для оплаты"
	ASKRECEIPTBUTTONTEXT        = "Отправка чека об оплате"
)

// After buttons texts
const (
	UNEXISTINGBUTTONPRESSED = "Неизвестная команда. Пожалуйста, используйте доступные кнопки."
	STARTGREETINGTEXT       = "Добро пожаловать в OneCoffee! Пожалуйста, отправьте ваше имя и номер телефона."
	POSTREGISTRATIONTEXT    = "Выберите действие1:"
	SENDCONTACTTEXT         = "Отправить cвой контакт"
	THANKYOUREGISTERTEXT    = "Спасибо, %s! Вы успешно зарегистрированы в OneCoffee."
	PAYMENTDETAILSTEXT      = "Вот реквизиты для оплаты:\n\nНазвание: OneCoffee\nНомер карты: 1234 5678 9012 3456\nБанк: ExampleBank\nНазначение платежа: подписка на кофе"
	ASKRECEIPTTEXT          = "Пожалуйста, отправьте фото или файл с чеком об оплате."
	MYSUBSCRIPTIONTEXT      = "Моя подписка"
)

// Error texts
const (
	ERRORREGISTRATIONTEXT = "Произошла ошибка при регистрации. Пожалуйста, попробуйте еще раз."
)
