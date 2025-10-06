import 'package:auto_light_pi/features/authentication/domain/use_cases/check_authentication.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/login.dart';
import 'package:auto_light_pi/features/authentication/domain/use_cases/logout.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_state.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_bloc.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_event.dart';
import 'package:auto_light_pi/features/login/presentation/bloc/login_state.dart';
import 'package:bloc_test/bloc_test.dart';
import 'package:mocktail/mocktail.dart';

/* Domain Mocks */
class MockLoginUseCase extends Mock implements LoginUseCase {}

class MockCheckAuthenticationUseCase extends Mock
    implements CheckAuthenticationUseCase {}

class MockLogoutUseCase extends Mock implements LogoutUseCase {}

/* Data Mocks */

/* Presentation Mocks */
class MockAuthenticationBloc
    extends MockBloc<AuthenticationEvent, AuthenticationState>
    implements AuthenticationBloc {}

class MockLoginBloc extends MockBloc<LoginEvent, LoginState>
    implements LoginBloc {}
