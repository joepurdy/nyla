/*
 * nyla-collector - GDPR compliant privacy focused web analytics
 * Copyright (C) 2024 Joe Purdy
 * mailto:nyla AT purdy DOT dev
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 3 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, write to the Free Software Foundation,
 * Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

interface CollectorData {
  type: "event" | "page";
  event: string;
  ua: string;
  hostname: string;
  referrer: string;
}

interface CollectPayload {
  data: CollectorData;
  site_id: string;
}

class Collector {
  private siteId: string = "";
  private referrer = "";

  constructor(siteId: string, ref: string) {
    this.siteId = siteId;
    this.referrer = ref;
  }

  track(event: string, type: "event" | "page") {
    const payload: CollectPayload = {
      data: {
        type: type,
        event: event,
        ua: navigator.userAgent,
        hostname: location.protocol + "//" + location.hostname,
        referrer: this.referrer
      },
      site_id: this.siteId
    };
    this.collectRequest(payload);
  }

  page(path: string) {
    this.track(path, "page");
  }

  private collectRequest(payload: CollectPayload) {
    const s = JSON.stringify(payload);
    const apiUrl = `http://localhost:9876/collect?data=${btoa(s)}`;

    const img = new Image();
    img.src = apiUrl;
  }
}

((w, d) => {
  const ds = d.currentScript?.dataset;
  if (!ds) {
    console.error("you must provide a data-siteid in your Nyla script tag.");
    return;
  } else if (!ds.siteid) {
    console.error("you must provide a data-siteid in your Nyla script tag.");
    return;
  }

  const path = w.location.pathname

  let externalReferrer = "";
  const ref = d.referrer;
  if (ref && ref.indexOf(`${w.location.protocol}//${w.location.host}`) == 0) {
    externalReferrer = ref;
  }

  let collector = new Collector(ds.siteid, externalReferrer);

  w._nyla = w._nyla || collector;

  collector.page(path);

  const his = window.history;
  if (his.pushState) {
    const originalFn = his["pushState"];
    his.pushState = function () {
      originalFn.apply(this, arguments);
      collector.page(w.location.pathname);
    };

    window.addEventListener("popstate", () => {
      collector.page(w.location.pathname);
    });
  }

  w.addEventListener(
    "hashchange",
    () => {
      collector.page(d.location.hash);
    },
    false
  );
})(window, document);