package sk.xco.eckcr.util;

import java.util.concurrent.TimeUnit;
import org.awaitility.Awaitility;
import org.awaitility.core.ThrowingRunnable;

public class Await {
  public static void untilAsserted(ThrowingRunnable runnable) {
    Awaitility.await().atMost(5, TimeUnit.SECONDS).untilAsserted(runnable);
  }
}
