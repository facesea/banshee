// Copyright 2015 Eleme Inc. All rights reserved.

const gulp = require('gulp');
const argv = require('yargs').argv;
const concat = require('gulp-concat');
const del = require('del');
const debug = require('gulp-debug');
const gulpif = require('gulp-if');
const minify = require('gulp-minify-css');
const uglify = require('gulp-uglify');

// Clean build
gulp.task('clean', function() {
  return del(['dist']);
});

// Css lib
gulp.task('css-lib', function() {
  var files = [
    'node_modules/bootstrap/dist/css/bootstrap.css',
  ];
  return gulp.src(files)
    .pipe(debug({title: 'Css-lib:'}))
    .pipe(minify())
    .pipe(concat('lib.min.css'))
    .pipe(debug({title: 'Css lib minified:'}))
    .pipe(gulp.dest('dist/css'));
});

// Css app
gulp.task('css-app', function() {
  var files = [
    'css/*.css',
  ];
  return gulp.src(files)
    .pipe(debug({title: 'Css:'}))
    .pipe(minify())
    .pipe(concat('app.min.css'))
    .pipe(debug({title: 'Css app minified:'}))
    .pipe(gulp.dest('dist/css'));
});

// Js lib
gulp.task('js-lib', function() {
  var files = [
    'node_modules/angular/angular.js',
    'node_modules/bootstrap/dist/js/bootstrap.js',
    'node_modules/d3/d3.js',
    'node_modules/cubism/cubism.v1.js',
  ];
  return gulp.src(files)
    .pipe(debug({title: 'Js-lib:'}))
    .pipe(gulpif(!argv.dev, uglify()))
    .pipe(concat('lib.min.js'))
    .pipe(debug({title: 'Js lib minified:'}))
    .pipe(gulp.dest('dist/js'));
});

// Js app
gulp.task('js-app', function() {
  var files = [
    'js/*.js',
  ];
  return gulp.src(files)
    .pipe(debug({title: 'Js-app:'}))
    .pipe(gulpif(!argv.dev, uglify()))
    .pipe(concat('app.min.js'))
    .pipe(debug({title: 'Js app minified:'}))
    .pipe(gulp.dest('dist/js'));
});

gulp.task('default', ['css-lib', 'css-app', 'js-lib', 'js-app']);
