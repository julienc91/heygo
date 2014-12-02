module.exports = function(grunt) {

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        nggettext_extract: {
            pot: {
                files: {
                    'translations/template.pot': ['static/html/admin/*.html', 'static/html/videos/*.html', 'static/js/*.js', 'templates/*.html']
                }
            },
        },
        nggettext_compile: {
            all: {
                files: {
                    'static/js/translations.js': ['translations/*.po']
                }
            },
        },
    })

    grunt.loadNpmTasks('grunt-angular-gettext');
    grunt.registerTask('default', ['nggettext_extract', 'nggettext_compile']);

}
