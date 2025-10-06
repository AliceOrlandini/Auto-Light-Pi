import 'package:equatable/equatable.dart';

abstract class LoginEvent extends Equatable {
  const LoginEvent();
  // coverage:ignore-start
  @override
  List<Object> get props => <Object>[];
  // coverage:ignore-end
}

final class LoginUsernameChanged extends LoginEvent {
  final String username;

  const LoginUsernameChanged(this.username);

  // This means: two instances of LoginUsernameChanged
  // are equal if their username properties are equal.
  // This is important for Bloc to distinguish new events
  // from old ones.
  @override
  List<Object> get props => <Object>[username];
}

final class LoginPasswordChanged extends LoginEvent {
  final String password;

  const LoginPasswordChanged(this.password);

  @override
  List<Object> get props => <Object>[password];
}

final class LoginSubmitted extends LoginEvent {
  // coverage:ignore-start
  const LoginSubmitted();
  // coverage:ignore-end
}
