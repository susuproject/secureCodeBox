/*
 *
 *  SecureCodeBox (SCB)
 *  Copyright 2015-2018 iteratec GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  	http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 * /
 */
package io.securecodebox.persistence.util;

import io.securecodebox.models.V1Scan;
import org.springframework.stereotype.Component;

import java.text.MessageFormat;
import java.time.Clock;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.Objects;

@Component
public class DescriptionGenerator {

  protected static final String TIME_FORMAT = "dd.MM.yyyy HH:mm:ss";
  private static final String DATE_FORMAT = "yyyy-MM-dd";
  Clock clock = Clock.systemDefaultZone();

  public String generate(V1Scan scan) {
    var spec = Objects.requireNonNull(scan.getSpec());

    return String.join(
      System.getProperty("line.separator"),
      MessageFormat.format("# {0}", getDefectDojoScanName(scan)),
      MessageFormat.format("Started: {0}", getStartTime(scan)),
      MessageFormat.format("Ended: {0}", getEndTime(scan)),
      MessageFormat.format("ScanType: {0}", spec.getScanType()),
      MessageFormat.format("Parameters: [{0}]", String.join(",", Objects.requireNonNull(spec.getParameters())))
    );
  }

  private String getStartTime(V1Scan scan) {
    if (scan.getMetadata() == null || scan.getMetadata().getCreationTimestamp() == null) {
      return null;
    }
    return scan.getMetadata().getCreationTimestamp().toString(TIME_FORMAT);
  }

  private String getEndTime(V1Scan scan) {
    if (scan.getStatus() == null || scan.getStatus().getFinishedAt() == null) {
      return currentTime();
    }
    return scan.getStatus().getFinishedAt().toString(TIME_FORMAT);
  }

  /**
   * Returns the current date as string based on the DATE_FORMAT.
   *
   * @return the current date as string based on the DATE_FORMAT.
   */
  public String currentDate() {
    return LocalDate.now(clock).format(DateTimeFormatter.ofPattern(DATE_FORMAT));
  }

  public String currentTime() {
    return LocalDateTime.now(clock).format(DateTimeFormatter.ofPattern(TIME_FORMAT));
  }

  public void setClock(Clock clock) {
    this.clock = clock;
  }

  public String getDefectDojoScanName(V1Scan scan) {
    return ScanNameMapping.bySecureCodeBoxScanType(scan.getSpec().getScanType()).scanType.getTestType();
  }
}
