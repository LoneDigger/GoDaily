const host = "http://localhost:8080";
axios.defaults.withCredentials = true;

function deleteRequset(url, params) {
  return axios.delete(host + url, params, {
    timeout: 30000,
  });
}

function getRequset(url, params) {
  return axios.get(host + url, params, {
    timeout: 30000,
  });
}

function postRequset(url, json) {
  return axios.post(host + url, json, {
    timeout: 30000,
    headers: {
      "Content-Type": "application/json",
    },
  });
}

function putRequset(url, json) {
  return axios.put(host + url, json, {
    timeout: 30000,
    headers: {
      "Content-Type": "application/json",
    },
  });
}

////////////////////////////////////////////////////////////////////////////////////////////////////

function hasClass(el, className) {
  if (el.classList)
    return el.classList.contains(className);
  return !!el.className.match(new RegExp('(\\s|^)' + className + '(\\s|$)'));
}

function addClass(el, className) {
  if (el.classList)
    el.classList.add(className)
  else if (!hasClass(el, className))
    el.className += " " + className;
}

function removeClass(el, className) {
  if (el.classList)
    el.classList.remove(className)
  else if (hasClass(el, className)) {
    var reg = new RegExp('(\\s|^)' + className + '(\\s|$)');
    el.className = el.className.replace(reg, ' ');
  }
}

////////////////////////////////////////////////////////////////////////////////////////////////////

// 月份
const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

// 每個月總和
function updateSumByMonth(sum, month) {
  const incomeChartConfig = {
    series: [
      {
        name: "",
        data: sum
      }
    ],
    chart: {
      parentHeightOffset: 0,
      parentWidthOffset: 0,
      toolbar: {
        show: false
      },
      type: 'area'
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 2,
      curve: 'smooth'
    },
    legend: {
      show: false
    },
    colors: [config.colors.primary],
    fill: {
      type: 'gradient',
      gradient: {
        //shade: shadeColor,
        shadeIntensity: 0.6,
        opacityFrom: 0.5,
        opacityTo: 0.25,
        stops: [0, 95, 100]
      }
    },
    grid: {
      borderColor: config.colors.borderColor,
      strokeDashArray: 3,
      padding: {
        top: -20,
        bottom: -8,
        left: 8,
        right: 8
      }
    },
    xaxis: {
      categories: month,
      axisBorder: {
        show: false
      },
      axisTicks: {
        show: false
      },
      labels: {
        show: true,
        style: {
          fontSize: '14px',
          colors: config.colors.axisColor
        }
      }
    },
    yaxis: {
      labels: {
        show: false
      },
      tickAmount: 4
    }
  };

  const incomeChart = document.querySelector("#incomeChart");
  const chart = new ApexCharts(incomeChart, incomeChartConfig);
  chart.render();
}

// 主類別分類總和
function updateSumByMainType(sum, name) {
  var total = 0;
  for (var i = 0; i < sum.length; i++) {
    total += sum[i];
  };

  const cardColor = config.colors.cardColor;
  const headingColor = config.colors.headingColor;
  const labelColor = config.colors.textMuted;
  const legendColor = config.colors.bodyColor;
  const borderColor = config.colors.borderColor;
  const axisColor = config.colors.axisColor;

  const doughnutChartConfig = {
    chart: {
      type: 'donut'
    },
    labels: name,
    series: sum,
    colors: ["#5fc9f8", "#fecb2e", "#fd9426", "#fc3158", "#147efb", "#53d769", "#fc3d39", "#8e8e93"],
    stroke: {
      colors: cardColor
    },
    dataLabels: {
      enabled: false,
      formatter: function (val, opt) {
        return parseInt(val);
      }
    },
    legend: {
      display: true,
      position: "top",
    },
    grid: {
      padding: {
        top: 0,
        bottom: 0,
        right: 15
      }
    },
    plotOptions: {
      pie: {
        donut: {
          size: "75%",
          labels: {
            show: true,
            value: {
              fontSize: "26px",
              fontFamily: "Public Sans",
              color: headingColor,
              offsetY: -18,
              formatter: function (val) {
                return parseInt(val);
              }
            },
            name: {
              offsetY: 20,
              fontFamily: "Public Sans"
            },
            total: {
              show: true,
              fontSize: "16px",
              color: axisColor,
              label: "花費",
              formatter: function (w) {
                return parseInt(total);
              }
            }
          }
        }
      }
    }
  };

  const doughnutChart = document.querySelector("#doughnutChart");
  const chart = new ApexCharts(doughnutChart, doughnutChartConfig);
  chart.render();
}
