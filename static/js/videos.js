$(function() {

    var slug = $("video").attr("data-slug");

    function resize_video() {
        $("video").width($(".video_container").width());
    }

    $(window).resize(resize_video);
    resize_video();

    var api_username = "";
    var api_password = "";
    var api_useragent = "OS Test User Agent";

    $("[role=search_subtitles]").click(function() {
        var data = {}

        get_hash(data);
        api_login(data);
        api_search_subtitles(data.hash, data.size, data.token);
        api_logout(data.token);
    });

    function api_login(d) {
        $.xmlrpc({

            url: 'http://api.opensubtitles.org/xml-rpc',
            methodName: 'LogIn',
            params: [api_username, api_password, "", api_useragent],
            async: false,

            success: function(response, status, jqXHR) {
                if (response.length > 0 && response[0].status == "200 OK") {
                    d.token = response[0].token;
                }
            }
        });
    }

    function api_logout(token) {
        $.xmlrpc({
            url: 'http://api.opensubtitles.org/xml-rpc',
            methodName: 'LogOut',
            params: [token]
        });
    }

    function api_search_subtitles(hash, size, token) {
        $.xmlrpc({

            url: 'http://api.opensubtitles.org/xml-rpc',
            methodName: 'SearchSubtitles',
            params: [token, [{"moviehash": hash, "moviebytesize": size}]],
            async: false,

            success: function(response, status, jqXHR) {
                console.log(response);
            }
        });
    }

    function get_hash(d) {
        $.ajax({

            url: "/videos/gethash/" + slug,
            async: false,

            success: function(data) {
                data = JSON.parse(data);
                if (data.ok) {
                    d.hash = data.hash;
                    d.size = data.size;
                }
            }
        });
    }
});
