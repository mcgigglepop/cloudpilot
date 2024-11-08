import { migratePatchFunc, migrateWarn } from '../main.js';

var oldRemoveAttr = jQuery.fn.removeAttr,
  oldToggleClass = jQuery.fn.toggleClass,
  rmatchNonSpace = /\S+/g;

migratePatchFunc(
  jQuery.fn,
  'removeAttr',
  function (name) {
    var self = this;

    jQuery.each(name.match(rmatchNonSpace), function (_i, attr) {
      if (jQuery.expr.match.bool.test(attr)) {
        migrateWarn(
          'removeAttr-bool',
          'jQuery.fn.removeAttr no longer sets boolean properties: ' + attr
        );
        self.prop(attr, false);
      }
    });

    return oldRemoveAttr.apply(this, arguments);
  },
  'removeAttr-bool'
);

migratePatchFunc(
  jQuery.fn,
  'toggleClass',
  function (state) {
    // Only deprecating no-args or single boolean arg
    if (state !== undefined && typeof state !== 'boolean') {
      return oldToggleClass.apply(this, arguments);
    }

    migrateWarn(
      'toggleClass-bool',
      'jQuery.fn.toggleClass( boolean ) is deprecated'
    );

    // Toggle entire class name of each element
    return this.each(function () {
      var className = (this.getAttribute && this.getAttribute('class')) || '';

      if (className) {
        jQuery.data(this, '__className__', className);
      }

      // If the element has a class name or if we're passed `false`,
      // then remove the whole classname (if there was one, the above saved it).
      // Otherwise bring back whatever was previously saved (if anything),
      // falling back to the empty string if nothing was stored.
      if (this.setAttribute) {
        this.setAttribute(
          'class',
          className || state === false
            ? ''
            : jQuery.data(this, '__className__') || ''
        );
      }
    });
  },
  'toggleClass-bool'
);
