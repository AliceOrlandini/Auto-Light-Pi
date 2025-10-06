import 'package:auto_light_pi/core/failure/network_failure.dart';
import 'package:auto_light_pi/features/authentication/data/models/login_valid_response.dart';
import 'package:auto_light_pi/network/dio_client.dart';
import 'package:dartz/dartz.dart';
import 'package:dio/dio.dart';

class AuthRemoteDataSource {
  final DioClient _dioClient;

  AuthRemoteDataSource(this._dioClient);

  Future<Either<LoginValidResponse, NetworkFailure>> loginByUsername(
    String username,
    String password,
  ) async {
    try {
      final Response<dynamic> response = await _dioClient.dio.post(
        '/login/username',
        data: <String, String>{'username': username, 'password': password},
      );
      return Left<LoginValidResponse, NetworkFailure>(
        LoginValidResponse.fromJson(response.data as Map<String, dynamic>),
      );
    } on DioException catch (e) {
      if (e.response != null && e.response!.data is Map<String, dynamic>) {
        final Map<String, dynamic> data =
            e.response!.data as Map<String, dynamic>;
        final String errorMsg = data['error'] as String? ?? 'unknown error';

        switch (e.response!.statusCode) {
          case 400:
            return Right<LoginValidResponse, NetworkFailure>(
              BadRequestFailure(errorMsg),
            );
          case 401:
            return Right<LoginValidResponse, NetworkFailure>(
              UnauthorizedFailure(errorMsg),
            );
          case 408:
            return Right<LoginValidResponse, NetworkFailure>(
              TimeoutFailure(errorMsg),
            );
          case 500:
            return Right<LoginValidResponse, NetworkFailure>(
              InternalServerFailure(errorMsg),
            );
          default:
            return Right<LoginValidResponse, NetworkFailure>(
              UnknownFailure(errorMsg),
            );
        }
      }
      // this is the Dio message, not the server response
      if (e.message != null) {
        return Right<LoginValidResponse, NetworkFailure>(
          UnknownFailure(e.message!),
        );
      }
      // if no response or exception was thrown, throw a generic exception
      return Right<LoginValidResponse, NetworkFailure>(
        UnknownFailure('Unknown error occurred'),
      );
    }
  }
}
