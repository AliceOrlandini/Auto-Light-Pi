import 'package:equatable/equatable.dart';
import 'package:formz/formz.dart';

enum PasswordValidationError {
  empty,
  tooShort,
  digitMissing,
  upperCaseMissing,
  lowerCaseMissing,
}

class Password extends FormzInput<String, PasswordValidationError>
    with EquatableMixin {
  const Password.pure() : super.pure('');
  const Password.dirty([super.value = '']) : super.dirty();

  @override
  List<Object> get props => <Object>[value];

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
    if (!value.contains(RegExp(r'[a-z]'))) {
      return PasswordValidationError.lowerCaseMissing;
    }
    return null;
  }
}
