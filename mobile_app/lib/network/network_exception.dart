class NetworkException implements Exception {
  final String message;
  final int? statusCode;

  NetworkException(this.message, {this.statusCode});

  @override
  String toString() => 'NetworkException($statusCode): $message';
}

class BadRequestException extends NetworkException {
  BadRequestException(super.msg) : super(statusCode: 400);
}

class UnauthorizedException extends NetworkException {
  UnauthorizedException(super.msg) : super(statusCode: 401);
}

class TimeoutException extends NetworkException {
  TimeoutException(super.msg) : super(statusCode: 408);
}

class ServerException extends NetworkException {
  ServerException(super.msg) : super(statusCode: 500);
}
