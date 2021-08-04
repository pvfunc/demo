'use strict';

const handler = async (event, context) => {
    await new Promise(resolve => setTimeout(resolve, 10000));

    const reply = "Long function executed!"

    return context
        .status(200)
        .succeed(reply);
};

module.exports = handler;
