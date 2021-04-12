/*
 * Copyright (c) 2020 Eric Kubenka.
 * Licensed under the MIT License.
 * See LICENSE file in the project root for full license information.
 */
package de.codefever.gladini;

import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.remote.RemoteWebDriver;
import org.testng.Assert;
import org.testng.annotations.Test;

import java.net.MalformedURLException;
import java.net.URL;

/**
 * @author Eric Kubenka
 * creation date: 16.06.2020 - 00:52
 */
public class SimpleSeleniumTest {

    private void nap() {
        try {
            long timeout = (long) (Math.random() * 10) * 1000;
            System.out.println("nap for " + timeout);
            Thread.sleep(timeout);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
    }

    @Test
    public void testT01_Start() throws MalformedURLException, InterruptedException {

        nap();
        final RemoteWebDriver driver = new RemoteWebDriver(new URL("http://localhost:30020"), new ChromeOptions());
        System.out.println("Session id is: " + driver.getSessionId());
        driver.get("https://google.de");

        for (int i = 0; i < 30; i++) {
            Thread.sleep(1000);
            Assert.assertEquals(driver.getTitle(), "Google");
        }

        driver.quit();
    }

    @Test
    public void testT02_Start() throws MalformedURLException {

        nap();
        final RemoteWebDriver driver = new RemoteWebDriver(new URL("http://localhost:30020"), new ChromeOptions());
        System.out.println("Session id is: " + driver.getSessionId());
        driver.get("https://google.de");
        nap();
        Assert.assertEquals(driver.getTitle(), "Google");
        driver.quit();
    }

    @Test
    public void testT03_Start() throws MalformedURLException {

        final RemoteWebDriver driver = new RemoteWebDriver(new URL("http://localhost:30020"), new ChromeOptions());
        System.out.println("Session id is: " + driver.getSessionId());
        driver.get("https://google.de");
        nap();
        Assert.assertEquals(driver.getTitle(), "Google");
        driver.quit();
    }

    @Test
    public void testT04_Start() throws MalformedURLException {

        final RemoteWebDriver driver = new RemoteWebDriver(new URL("http://localhost:30020"), new ChromeOptions());
        System.out.println("Session id is: " + driver.getSessionId());
        driver.get("https://google.de");
        nap();
        Assert.assertEquals(driver.getTitle(), "Google");
        driver.quit();
    }

    @Test
    public void testT05_Start() throws MalformedURLException {

        nap();
        final RemoteWebDriver driver = new RemoteWebDriver(new URL("http://localhost:30020"), new ChromeOptions());
        System.out.println("Session id is: " + driver.getSessionId());
        driver.get("https://google.de");
        nap();
        Assert.assertEquals(driver.getTitle(), "Google");
        driver.quit();
    }


}
