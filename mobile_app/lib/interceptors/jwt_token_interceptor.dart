import 'package:auto_light_pi/storage/jwt_token_storage.dart';
import 'package:dio/dio.dart';

class JwtTokenInterceptor extends Interceptor {
  final JwtTokenStorage _jwtTokenStorage;

  JwtTokenInterceptor(this._jwtTokenStorage);

  @override
  void onResponse(
    Response<dynamic> response,
    ResponseInterceptorHandler handler,
  ) async {
    final List<String>? setCookie = response.headers.map['set-cookie'];
    if (setCookie != null) {
      final RegExpMatch? match = RegExp(
        r'jwt=([^;]+)',
      ).firstMatch(setCookie.join(';'));
      if (match != null) {
        await _jwtTokenStorage.writeToken(match.group(1)!);
      }
    }
    handler.next(response);
  }

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final String? token = await _jwtTokenStorage.readToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }
}
