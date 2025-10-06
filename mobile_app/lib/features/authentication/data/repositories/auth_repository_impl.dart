import 'package:auto_light_pi/core/failure/network_failure.dart';
import 'package:auto_light_pi/features/authentication/data/data_sources/auth_remote_data_source.dart';
import 'package:auto_light_pi/features/authentication/data/models/login_valid_response.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/domain/repositories/auth_repository.dart';
import 'package:auto_light_pi/features/authentication/data/data_sources/auth_local_data_source.dart';
import 'package:dartz/dartz.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource _remote;
  final AuthLocalDataSource _local;

  AuthRepositoryImpl({
    required AuthRemoteDataSource remote,
    required AuthLocalDataSource local,
  }) : _remote = remote,
       _local = local;

  @override
  Future<Either<UserEntity, NetworkFailure>> authenticate({
    required String username,
    required String password,
  }) async {
    final Either<LoginValidResponse, NetworkFailure> result = await _remote
        .loginByUsername(username, password);
    return result.fold(
      (LoginValidResponse loginValidResponse) {
        final UserEntity user = loginValidResponse.user.toEntity();

        // Cache the user locally
        _local.cacheUser(loginValidResponse.user);
        return Left<UserEntity, NetworkFailure>(user);
      },
      (NetworkFailure failure) {
        return Right<UserEntity, NetworkFailure>(failure);
      },
    );
  }

  @override
  Future<UserEntity?> isAuthenticated() async {
    final UserEntity? user = await _local.getUser();
    return user;
  }

  @override
  Future<void> logout() async {
    await _local.clearUser();
    await _local.clearToken();
  }
}
