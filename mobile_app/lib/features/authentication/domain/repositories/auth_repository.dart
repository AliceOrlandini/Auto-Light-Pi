import 'package:dartz/dartz.dart';
import 'package:auto_light_pi/core/failure/network_failure.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';

/// Try to authenticate the user using username and password and it returns
/// an UserEntity in case of success, a Failure otherwise
abstract class AuthRepository {
  Future<Either<UserEntity, NetworkFailure>> authenticate({
    required String username,
    required String password,
  });

  /// Check if the user is authenticated
  Future<UserEntity?> isAuthenticated();

  /// Logout the user by clearing the cached data
  Future<void> logout();
}
