module.exports = function(grunt) {

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        nggettext_extract: {
            pot: {
                files: {
                    'translations/template.pot': ['static/html/*/*.html', 'templates/*.html']
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
