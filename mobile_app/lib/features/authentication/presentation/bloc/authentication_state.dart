import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:equatable/equatable.dart';

enum AuthenticationStatus { unknown, authenticated, unauthenticated }

final class AuthenticationState extends Equatable {
  final AuthenticationStatus status;
  final UserEntity user;
  final String? errorMessage;
  final int? statusCode;

  const AuthenticationState._({
    this.status = AuthenticationStatus.unknown,
    this.user = UserEntity.empty,
    this.errorMessage,
    this.statusCode,
  });

  const AuthenticationState.unknown() : this._();

  const AuthenticationState.authenticated(UserEntity user)
    : this._(status: AuthenticationStatus.authenticated, user: user);

  const AuthenticationState.unauthenticated(
    String? errorMessage, {
    int? statusCode,
  }) : this._(
         status: AuthenticationStatus.unauthenticated,
         errorMessage: errorMessage,
         statusCode: statusCode,
       );

  @override
  List<Object> get props => <Object>[status, user];
}
