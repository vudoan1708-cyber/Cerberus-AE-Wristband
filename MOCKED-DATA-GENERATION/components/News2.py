def Generate(respiration, sp02, airNotOxygen, pulse, temperature, bloodPressure, motion, overall):
  
    return dict([
      ("respiration", respiration),
      ("sp02", sp02),
      ("onOxygen", airNotOxygen),
      ("pulse", pulse),
      ("temperature", temperature),
      ("bloodPressure", bloodPressure),
      ("motion", motion),
      # ("proximity", proximity),
      ("overall", overall),
    ])