plugins {
    alias(libs.plugins.android.library) apply false
    alias(libs.plugins.kotlin.android) apply false
    alias(libs.plugins.jetbrains.kotlin.jvm)
    `java-library`
}

allprojects {
    tasks.withType<JavaCompile> {
        options.encoding = "UTF-8"
    }
}
