var app = angular.module('gomet', ['pascalprecht.translate'], function($interpolateProvider) {
    $interpolateProvider.startSymbol('[[');
    $interpolateProvider.endSymbol(']]');
});

app.config(['$translateProvider', function ($translateProvider) {
    $translateProvider.translations('fr_FR', {
        'ADMIN_TITLE': 'Administration',
        'ADMIN_USERS_TAB': 'Utilisateurs',
        'ADMIN_INVITATIONS_TAB': 'Invitations',
        'ADMIN_GROUPS_TAB': 'Groupes d\'utilisateurs',
        'ADMIN_VIDEOS_TAB': 'Vidéos',
        'ADMIN_VIDEO_GROUPS_TAB': 'Groupes de vidéos',
        'ADMIN_PERMISSIONS_TAB': 'Permissions',
        'ADMIN_MEMBERSHIP_TAB': 'Groupes d\'utilisateurs',
        'ADMIN_CLASSIFICATION_TAB': 'Catégories de vidéos',
        
        'ADD': 'Ajouter',
        'CANCEL': 'Annuler',
        'CLOSE': 'Fermer',
        'DELETE': 'Supprimer',
        'ID': 'Id',
        'LOGIN': 'Login',
        'NAME': 'Nom',
        'PASSWORD': 'Mot de passe',
        'PATH': 'Chemin',
        'RESET': 'Réinitialiser',
        'SAVE': 'Enregistrer',
        'SLUG': 'Slug',
        'TITLE': 'Nom',
        'URL': 'Url',
        'VALUE': 'Valeur',
        
        'MENU_ABOUT': 'A propos',
        'MENU_VIDEOS': 'Vidéos',
        'MENU_ADMINISTRATION': 'Administration',
        'MENU_SIGNOUT': 'Déconnexion'
    });
    
    $translateProvider.preferredLanguage('fr_FR');
}]);
