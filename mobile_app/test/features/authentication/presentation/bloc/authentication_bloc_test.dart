import 'package:auto_light_pi/core/failure/network_failure.dart';
import 'package:auto_light_pi/features/authentication/domain/entities/user_entity.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:dartz/dartz.dart';
import 'package:mocktail/mocktail.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:flutter_test/flutter_test.dart';

import '../../../../mocks/mocks.dart';

void main() {
  late MockLoginUseCase mockLoginUseCase;
  late MockCheckAuthenticationUseCase mockCheckAuthenticationUseCase;
  late MockLogoutUseCase mockLogoutUseCase;

  setUp(() {
    mockLoginUseCase = MockLoginUseCase();
    mockCheckAuthenticationUseCase = MockCheckAuthenticationUseCase();
    mockLogoutUseCase = MockLogoutUseCase();
  });

  group('AuthenticationBloc', () {
    test(
      'initial state is AuthenticationState.unknown, User is empty, errorMessage and statusCode are null',
      () {
        final AuthenticationBloc authenticationBloc = AuthenticationBloc(
          loginUseCase: mockLoginUseCase,
          checkAuthenticationUseCase: mockCheckAuthenticationUseCase,
          logoutUseCase: mockLogoutUseCase,
        );
        expect(authenticationBloc.state, const AuthenticationState.unknown());
        expect(authenticationBloc.state.user, UserEntity.empty);
        expect(authenticationBloc.state.errorMessage, isNull);
        expect(authenticationBloc.state.statusCode, isNull);
        authenticationBloc.close();
      },
    );

    blocTest<AuthenticationBloc, AuthenticationState>(
      'emits [AuthenticationState.unknown(), AuthenticationState.authenticated] when LoginRequest is added and login is successful',
      build: () {
        when(
          () => mockLoginUseCase.call(
            username: any(named: 'username'),
            password: any(named: 'password'),
          ),
        ).thenAnswer(
          (_) async => const Left<UserEntity, NetworkFailure>(
            UserEntity(
              username: 'testuser',
              email: 'testuser@mail.com',
              name: 'Test',
              surname: 'User',
            ),
          ),
        );
        return AuthenticationBloc(
          loginUseCase: mockLoginUseCase,
          checkAuthenticationUseCase: mockCheckAuthenticationUseCase,
          logoutUseCase: mockLogoutUseCase,
        );
      },
      act: (AuthenticationBloc bloc) =>
          bloc.add(const LoginRequest('testuser', 'password123')),
      expect: () => <AuthenticationState>[
        const AuthenticationState.unknown(),
        const AuthenticationState.authenticated(
          UserEntity(
            username: 'testuser',
            email: 'testuser@mail.com',
            name: 'Test',
            surname: 'User',
          ),
        ),
      ],
    );

    blocTest<AuthenticationBloc, AuthenticationState>(
      'emits [AuthenticationState.unknown(), AuthenticationState.unauthenticated] when LoginRequest is added and login fails',
      build: () {
        when(
          () => mockLoginUseCase.call(
            username: any(named: 'username'),
            password: any(named: 'password'),
          ),
        ).thenAnswer(
          (_) async => Right<UserEntity, NetworkFailure>(
            BadRequestFailure('Invalid credentials'),
          ),
        );
        return AuthenticationBloc(
          loginUseCase: mockLoginUseCase,
          checkAuthenticationUseCase: mockCheckAuthenticationUseCase,
          logoutUseCase: mockLogoutUseCase,
        );
      },
      act: (AuthenticationBloc bloc) =>
          bloc.add(const LoginRequest('testuser', 'wrongpassword')),
      expect: () => <AuthenticationState>[
        const AuthenticationState.unknown(),
        const AuthenticationState.unauthenticated(
          'Invalid credentials',
          statusCode: 401,
        ),
      ],
    );

    blocTest<AuthenticationBloc, AuthenticationState>(
      'emits [AuthenticationState.unknown(), AuthenticationState.authenticated] when AuthCheckRequest is added and user is authenticated',
      build: () {
        when(() => mockCheckAuthenticationUseCase.call()).thenAnswer(
          (_) async => const UserEntity(
            username: 'testuser',
            email: 'testuser@mail.com',
            name: 'Test',
            surname: 'User',
          ),
        );
        return AuthenticationBloc(
          loginUseCase: mockLoginUseCase,
          checkAuthenticationUseCase: mockCheckAuthenticationUseCase,
          logoutUseCase: mockLogoutUseCase,
        );
      },
      act: (AuthenticationBloc bloc) => bloc.add(AuthCheckRequest()),
      expect: () => <AuthenticationState>[
        const AuthenticationState.unknown(),
        const AuthenticationState.authenticated(
          UserEntity(
            username: 'testuser',
            email: 'testuser@mail.com',
            name: 'Test',
            surname: 'User',
          ),
        ),
      ],
    );

    blocTest<AuthenticationBloc, AuthenticationState>(
      'emits [AuthenticationState.unknown(), AuthenticationState.unauthenticated] when AuthCheckRequest is added and user is unauthenticated',
      build: () {
        when(
          () => mockCheckAuthenticationUseCase.call(),
        ).thenAnswer((_) async => null);
        return AuthenticationBloc(
          loginUseCase: mockLoginUseCase,
          checkAuthenticationUseCase: mockCheckAuthenticationUseCase,
          logoutUseCase: mockLogoutUseCase,
        );
      },
      act: (AuthenticationBloc bloc) => bloc.add(AuthCheckRequest()),
      expect: () => <AuthenticationState>[
        const AuthenticationState.unknown(),
        const AuthenticationState.unauthenticated('Unauthenticated'),
      ],
    );

    blocTest<AuthenticationBloc, AuthenticationState>(
      'emits [AuthenticationState.unauthenticated] when LogoutRequest is added',
      build: () {
        when(
          () => mockLogoutUseCase.call(),
        ).thenAnswer((_) async => <dynamic, dynamic>{});
        return AuthenticationBloc(
          loginUseCase: mockLoginUseCase,
          checkAuthenticationUseCase: mockCheckAuthenticationUseCase,
          logoutUseCase: mockLogoutUseCase,
        );
      },
      act: (AuthenticationBloc bloc) => bloc.add(LogoutRequest()),
      expect: () => <AuthenticationState>[
        const AuthenticationState.unauthenticated('Logged out'),
      ],
    );
  });
}
