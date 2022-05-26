package handlers;

import ru.sber.platformv.faas.api.HttpFunction;
import ru.sber.platformv.faas.api.HttpRequest;
import ru.sber.platformv.faas.api.HttpResponse;

import java.io.IOException;
import java.util.logging.Logger;

public class Handler2 implements HttpFunction {

    // Метод service. Данный метод будет обрабатывать HTTP запросы поступающие к функции
    @Override
    public void service(HttpRequest request, HttpResponse response) throws IOException {

        // Логирование входящего запроса
        String requestBody = new String(request.getInputStream().readAllBytes());
        Logger.getGlobal().info("Request received: " + requestBody + "\nMethod: " + request.getMethod());

        // Подготовка и возврат ответа на вызов
        response.setContentType("text/plain; charset=utf-8");
        response.getWriter().write("Hello from Java11 function!\nYou said: " + requestBody);
    }
}