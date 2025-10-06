import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/domain/repositories/auth_repository.dart';

class CheckAuthenticationUseCase {
  final AuthRepository _authRepository;

  CheckAuthenticationUseCase(this._authRepository);

  Future<UserEntity?> call() {
    return _authRepository.isAuthenticated();
  }
}
