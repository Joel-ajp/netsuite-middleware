/**
 * @NApiVersion 2.x
 * @NScriptType UserEventScript
 */
define(['N/https', 'N/log', 'N/runtime', 'N/record'], function(https, log, runtime) {
    function get_test(context) {
        var newRecord = context.newRecord;
        var scriptObj = runtime.getCurrentScript();
        var get_endpoint = scriptObj.getParameter({name: "custscript1"});
        var post_endpoint = scriptObj.getParameter({name: "custscript2"});

        try {
            var client = newRecord.getValue("inpt_parent");
            var project = newRecord.getValue("entityid");
            var subsidiary = newRecord.getValue("inpt_subsidiary");
            var form_type = newRecord.getValue("inpt_customform");
        } catch (e) {
            log.error({
                title: 'Error',
                details: e
            });
            return
        }

        if (client == "" || project == "") {
            log.debug({
                title: 'Error',
                details: 'Client or Project is empty'
            });
            return;
        }

        var json = {
            "client": client,
            "project": project,
            "subsidiary": subsidiary,
            "form_type": form_type
        }

        var response = https.post({
            url: post_endpoint,
            body: JSON.stringify(json),
            headers: {
                'Content-Type': 'application/json'
            }
        });

        log.debug({
            title: 'API Response',
            details: response.body
        });
    }

    return {
        afterSubmit: get_test,
    };
});
