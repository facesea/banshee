var pkg = require('./package.json'),
  gulp = require('gulp'),
  concat = require('gulp-concat'),
  uglify = require('gulp-uglify'),
  rename = require('gulp-rename'),
  htmlmin = require('gulp-htmlmin'),
  templateCache = require('gulp-angular-templatecache'),
  usemin = require('gulp-usemin'),
  minifyCss = require('gulp-minify-css'),
  rev = require('gulp-rev'),
  connect = require('gulp-connect'),
  gulpSequence = require('gulp-sequence'),
  browserify = require('gulp-browserify'),
  less = require('gulp-less');

var url = require('url');
var proxy = require('proxy-middleware');
var modRewrite = require('connect-modrewrite');

var jshint = require('gulp-jshint');
var stylish = require('jshint-stylish');

var dist = './public/js';

gulp.task('default', ['dev']);

// Less
gulp.task('less', function () {
  return gulp.src(['./src/theme/inspinia.less', './src/app.less', './src/theme/square/green.css'])
    .pipe(less())
    .pipe(gulp.dest('./public/css'))
    .pipe(connect.reload());
});

// Js
gulp.task('js', function () {
  var path;
  if (this.seq.indexOf('build') == -1) {
    path = 'public';
  } else {
    path = 'dist';
  };

  gulp.src('src/app.js')
    .pipe(browserify({
      debug: true,
      insertGlobals: true,
      transform: ['browserify-ngannotate']
    }))
    .pipe(rename(pkg.name + '.js'))
    .pipe(gulp.dest(dist))
    .pipe(connect.reload());
});

// Lint
gulp.task('lint', function () {
  return gulp.src([
      'src/**/*.js',
      '!src/mock.js',
      '!src/mock/*.js'
    ])
    .pipe(jshint())
    .pipe(jshint.reporter(stylish))
    .pipe(jshint.reporter('gulp-jshint-html-reporter', {
      filename: __dirname + '/jshint-output.html',
      createMissingFolders: false
    }));
});

// Html
gulp.task('tpl', function () {
  return gulp.src(['src/**/**.html'])
    .pipe(htmlmin({
      collapseBooleanAttributes: true,
      collapseWhitespace: true,
      removeAttributeQuotes: true,
      removeComments: true,
      removeEmptyAttributes: true,
      removeRedundantAttributes: true,
      removeScriptTypeAttributes: true,
      removeStyleLinkTypeAttributes: true
    }))
    .pipe(templateCache({
      filename: pkg.name + '.tpl.js',
      module: pkg.name + '.tpl',
      standalone: true
    }))
    .pipe(gulp.dest(dist))
    .pipe(connect.reload());
});

// File
gulp.task('file', function () {
  var path;
  if (this.seq.indexOf('build') == -1) {
    path = 'public';
  } else {
    path = 'dist';
  };

  gulp.src(['./public/node_modules/bootstrap/fonts/*', './public/node_modules/font-awesome/fonts/*'])
    .pipe(gulp.dest('./' + path + '/fonts'));

  gulp.src(['./public/images/**/*', './src/theme/*.png', './src/theme/square/*.png', './src/static/*.png', './src/static/*.ico'])
    .pipe(gulp.dest('./' + path + '/images'));
});

// Minify
gulp.task('usemin', ['tpl', 'less', 'js'], function () {
  return gulp.src('./public/index.html')
    .pipe(usemin({
      cssLib: [minifyCss(), 'concat', rev()],
      cssApp: [minifyCss(), 'concat', rev()],
      jsLib: ['concat', rename({
        suffix: '.min'
      }), rev()],
      jsApp: ['concat', uglify(), rename({
        suffix: '.min'
      }), rev()]
    }))
    .pipe(gulp.dest('dist/'));
});

// Build
gulp.task('build', gulpSequence('usemin', 'file'));

// Dev
gulp.task('dev', ['lint', 'tpl', 'file', 'less', 'js'], function () {
  // Server
  connect.server({
    root: 'public',
    port: 3000,
    livereload: true,
    middleware: function (connect, opt) {
      var middlewares = [];
      try {
        middlewares = require('./proxy.json').map(function (opt) {
          return proxy(opt);
        });
      } catch (e) {
        console.warn('proxy.json not found');
      }
      middlewares.push(modRewrite(['^[^\\.]*$ /index.html [L]']));
      return middlewares;
    }
  });
  gulp.watch('src/**/**.html', ['tpl']);
  gulp.watch('src/**/*.js', ['js']);
  gulp.watch('src/**/*.less', ['less']);
});
