'use strict';

const handler = async (event, context) => {
    const request = "Request: " + JSON.stringify(event);
    console.log(request)

    const reply = "Callback function.\n" +
        `${request}`

    return context
        .status(200)
        .succeed(reply);
};

module.exports = handler;
