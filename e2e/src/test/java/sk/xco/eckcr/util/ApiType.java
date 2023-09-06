package sk.xco.eckcr.util;

import lombok.Getter;

@Getter
public enum ApiType {
  Index("index.es.eck.github.com"),
  IndexTemplate("indextemplate.es.eck.github.com"),
  IndexLifecyclePolicy("indexlifecyclepolicy.es.eck.github.com");

  private final String resourceType;

  ApiType(String resourceType) {
    this.resourceType = resourceType;
  }
}
