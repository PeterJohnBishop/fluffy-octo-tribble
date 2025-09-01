import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

import 'main_views/main_carousel.dart';
import 'main_views/main_eddie.dart';
import 'main_views/main_nav.dart';
import 'main_views/main_tyler.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {

    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.grey),
        scaffoldBackgroundColor: Colors.white,
        textTheme: GoogleFonts.alumniSansPinstripeTextTheme(),
      ),
      home: MainView(),
    );
  }
}

class MainView extends StatefulWidget {
  const MainView({super.key});

  @override
  State<MainView> createState() => _MainViewState();
}

class _MainViewState extends State<MainView> {
  late int selectedIndex;

  @override
  void initState() {
    selectedIndex = 0;
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    var height = MediaQuery.of(context).size.height;
    var width = MediaQuery.of(context).size.width;

    Widget getContentView() {
      switch (selectedIndex) {
        case 0:
          return MainCarouselView();
        case 1:
          return SingleChildScrollView(
    child: Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: const [
        EddieView(),           // sizes itself based on content
        SizedBox(height: 50),  // spacing
        TylerView(),           // sizes itself based on content
        SizedBox(height: 50),
      ],
    ),
  );
        case 2:
          return Center(child: Text('Listings View', style: TextStyle(fontSize: 30)));
        case 3:
          return Center(child: Text('Buyers View', style: TextStyle(fontSize: 30)));
        case 4:
          return Center(child: Text('Sellers View', style: TextStyle(fontSize: 30)));
        case 5:
          return Center(child: Text('Contact View', style: TextStyle(fontSize: 30)));
        default:
          return MainCarouselView();
      }
    }

    return Scaffold(
      body: Column(
        children: [
          Container(
            width: double.infinity,
            color: Colors.white,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              mainAxisSize: MainAxisSize.min,
              children: [
                Text(
                  'LEDERMAN LOCATIONS',
                  style: const TextStyle(
                    fontSize: 42,
                    color: Colors.black,
                  ),
                ),
                Text(
                  'Residential & Commercial Real Estate',
                  style: TextStyle(
                    fontSize: 24,
                    fontWeight: FontWeight.w200,
                    color: Colors.blueGrey,
                  ),
                ),
              ],
            ),
          ),
          // Navigation
          MainNavView(
            height: height,
            width: width,
            selectedIndex: selectedIndex,
            onIndexChanged: (index) {
              setState(() {
                selectedIndex = index;
              });
            },
          ),
          // Main content fills remaining space
          Expanded(child: getContentView()),
        ],
      ),
    );
  }
}
