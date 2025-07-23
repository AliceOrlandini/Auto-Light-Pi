import 'package:auto_light_pi/core/errors/backend_failures.dart';
import 'package:auto_light_pi/core/errors/failure.dart';
import 'package:auto_light_pi/features/authentication/data/data_sources/auth_remote_data_source.dart';
import 'package:auto_light_pi/features/authentication/data/models/login_response.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/domain/repositories/auth_repository.dart';
import 'package:auto_light_pi/network/network_exception.dart';
import 'package:dartz/dartz.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthRemoteDataSource _remote;
  //final AuthLocalDataSource _local;

  AuthRepositoryImpl({required AuthRemoteDataSource remote}) : _remote = remote;

  @override
  Future<Either<UserEntity, Failure>> authenticate({
    required String username,
    required String password,
  }) async {
    try {
      final LoginResponse loginResponse = await _remote.loginByUsername(
        username,
        password,
      );
      final UserEntity user = loginResponse.user.toEntity();
      return Left<UserEntity, Failure>(user);
    } on BadRequestException catch (e) {
      return Right<UserEntity, Failure>(BadRequestFailure(e.message));
    } on UnauthorizedException catch (e) {
      return Right<UserEntity, Failure>(UnauthorizedFailure(e.message));
    } on TimeoutException catch (e) {
      return Right<UserEntity, Failure>(TimeoutFailure(e.message));
    } on ServerException catch (e) {
      return Right<UserEntity, Failure>(ServerFailure(e.message));
    } on NetworkException catch (e) {
      return Right<UserEntity, Failure>(UnknownFailure(e.message));
    } catch (e) {
      return Right<UserEntity, Failure>(UnknownFailure(e.toString()));
    }
  }
}
