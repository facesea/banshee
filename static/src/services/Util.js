/*@ngInject*/
module.exports = function () {
  return {
    buildRepr: function (rule) {
      var parts = [];

      var trendUp = rule.trendUp || false;
      var trendDown = rule.trendDown || false;
      var thresholdMax = rule.thresholdMax || 0;
      var thresholdMin = rule.thresholdMin || 0;

      if (trendUp && thresholdMax === 0) {
        parts.push('trend ↑');
      }
      if (trendUp && thresholdMax !== 0) {
        parts.push('(trend ↑ && value >= ' + parseFloat(thresholdMax.toFixed(3)) + ')');
      }
      if (!trendUp && thresholdMax !== 0) {
        parts.push('value >= ' + parseFloat(thresholdMax.toFixed(3)));
      }
      if (trendDown && thresholdMin === 0) {
        parts.push('trend ↓');
      }
      if (trendDown && thresholdMin !== 0) {
        parts.push('(trend ↓ && value <= ' + parseFloat(thresholdMin.toFixed(3)) + ')');
      }
      if (!trendDown && thresholdMin !== 0) {
        parts.push('value <= ' + parseFloat(thresholdMin.toFixed(3)));
      }
      return parts.join(' || ');
    }
  };
};
