import 'package:auto_light_pi/di/di.dart' as di;
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_bloc.dart';
import 'package:auto_light_pi/features/authentication/presentation/bloc/authentication_event.dart';
import 'package:flutter/material.dart';
import 'package:auto_light_pi/navigation/routes.dart';
import 'package:auto_light_pi/core/theme/app_theme.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:provider/single_child_widget.dart';

class App extends StatelessWidget {
  const App({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiBlocProvider(
      providers: <SingleChildWidget>[
        BlocProvider<AuthenticationBloc>(
          create: (_) => di.sl<AuthenticationBloc>()..add(AuthCheckRequest()),
        ),
      ],
      child: MaterialApp.router(
        title: 'Auto Light Pi',
        routerConfig: router,
        debugShowCheckedModeBanner: false,
        theme: AppTheme.getTheme(ThemeMode.light),
        darkTheme: AppTheme.getTheme(ThemeMode.dark),
        themeMode: ThemeMode.dark,
      ),
    );
  }
}
