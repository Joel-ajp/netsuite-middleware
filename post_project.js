/**
 * @NApiVersion 2.x
 * @NScriptType UserEventScript
 */
define(['N/https', 'N/log', 'N/runtime', 'N/record'], function(https, log, runtime) {
    function get_test(context) {
        var newRecord = context.newRecord;
        var scriptObj = runtime.getCurrentScript();
        var get_test = scriptObj.getParameter({name: "custscript1"});

        var post_endpoint = scriptObj.getParameter({name: 'post_id'});

        var response = https.get({
            url: get_test,
        });

        log.debug({
            title: 'API Response',
            details: response.body
        });

        var fields = newRecord.getFields();

        for (var i = 0; i < fields.length; i++) {
            var field = fields[i];
            var value = newRecord.getValue(field);
            log.debug('Field:' + field, 'Value:' + value);
        }
    }

    return {
        afterSubmit: get_test,
    };
});
