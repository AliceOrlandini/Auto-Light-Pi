import 'package:auto_light_pi/core/failure/network_failure.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/domain/repositories/auth_repository.dart';
import 'package:dartz/dartz.dart';

class LoginUseCase {
  final AuthRepository _authRepository;

  LoginUseCase(this._authRepository);

  Future<Either<UserEntity, NetworkFailure>> call({
    required String username,
    required String password,
  }) {
    return _authRepository.authenticate(username: username, password: password);
  }
}
