import 'package:formz/formz.dart';

enum PasswordValidationError { empty, tooShort, digitMissing, upperCaseMissing }

class Password extends FormzInput<String, PasswordValidationError> {
  const Password.pure() : super.pure('');
  const Password.dirty([super.value = '']) : super.dirty();

  @override
  PasswordValidationError? validator(String value) {
    if (value.isEmpty) {
      return PasswordValidationError.empty;
    }
    if (value.length < 8) {
      return PasswordValidationError.tooShort;
    }
    if (!value.contains(RegExp(r'[0-9]'))) {
      return PasswordValidationError.digitMissing;
    }
    if (!value.contains(RegExp(r'[A-Z]'))) {
      return PasswordValidationError.upperCaseMissing;
    }
    return null;
  }
}
